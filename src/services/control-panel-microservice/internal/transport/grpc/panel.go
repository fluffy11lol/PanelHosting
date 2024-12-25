package transport

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	interceptors "control-panel/internal/interceptors"
	models "control-panel/internal/models"
	database "control-panel/internal/repository"
	panel "control-panel/pkg/api/panel"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var DB *database.PanelStorage

type controlPanel struct {
	panel.UnimplementedPanelServiceServer
}

func NewControlPanel(ctx context.Context, host_psql, user_psql, password_psql, dbname_psql, port_psql, host_mysql, user_mysql, password_mysql, port_mysql, dbname_mysql string) (*controlPanel, error) {
	//dsnPSQL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host_psql, user_psql, password_psql, dbname_psql, port_psql)
	//dsnMySQL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	//	user_mysql, password_mysql, host_mysql, port_mysql, dbname_mysql)
	//
	//db, err := DB.StartDB(ctx, dsnPSQL, dsnMySQL)
	//if err != nil {
	//	return nil, err
	//}
	//DB = db
	//err = DB.DB.AutoMigrate(&models.User{})
	//if err != nil {
	//	log.Fatalf("failed migration: %v", err)
	//}
	//err = DB.DB.AutoMigrate(&models.Server{})
	//if err != nil {
	//	log.Fatalf("failed migration: %v", err)
	//}
	//err = DB.DB.AutoMigrate(&models.ServerDatabases{})
	//if err != nil {
	//	log.Fatalf("failed migration: %v", err)
	//}
	//
	//fmt.Printf("Serving PostgreSQL on 0.0.0.0:%s\n", port_psql)
	//fmt.Printf("Serving MySQL on 0.0.0.0:%s\n", port_mysql)

	return &controlPanel{}, nil
}

func (c *controlPanel) Login(ctx context.Context, req *panel.LoginRequest) (*panel.LoginResponse, error) {
	id, storedPassword, err := DB.GetUserByName(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Ошибка, если пользователь не найден
			return nil, status.Errorf(codes.NotFound, "user %s not found", req.Username)
		}
		// Другая ошибка
		return nil, status.Errorf(codes.Internal, "error fetching user %s: %v", req.Username, err)
	}

	if storedPassword != req.Password {
		return nil, status.Errorf(codes.Unauthenticated, fmt.Sprintf("invalid password for user %s", req.Username))
	}

	token, err := interceptors.GenerateJWTtoken(id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to generate JWT token for user %s", req.Username))
	}

	return &panel.LoginResponse{Token: token}, nil
}

func (c *controlPanel) CreateServer(ctx context.Context, req *panel.CreateServerRequest) (*panel.CreateServerResponse, error) {
	ID, err := strconv.Atoi(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid user ID: %v", err))
	}

	server, err := DB.CreateServer(req.Name, "", ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to create server: %v", err))
	}

	Server := panel.Server{
		Id:      string(server.ID),
		Name:    server.Name,
		Address: server.Address,
		Port:    server.Port,
		Status:  server.Status,
	}

	return &panel.CreateServerResponse{Server: &Server}, nil
}

func (c *controlPanel) GetServerDetails(ctx context.Context, req *panel.ServerDetailsRequest) (*panel.ServerDetailsResponse, error) {
	server, err := DB.GetServerDetails(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("server not found: %v", err))
	}

	serverDetails := panel.Server{
		Id:           server.ID,
		Name:         server.Name,
		Address:      server.Address,
		Port:         server.Port,
		Status:       server.Status,
		Tariffstatus: server.TariffStatus,
	}

	return &panel.ServerDetailsResponse{Server: &serverDetails}, nil
}

func (c *controlPanel) ListServers(ctx context.Context, req *panel.ListServersRequest) (*panel.ListServersResponse, error) {
	if req.UserId == "" {
		log.Println("user id is empty")
		return nil, status.Errorf(codes.InvalidArgument, "user ID is empty")
	}

	ID, err := strconv.Atoi(req.UserId)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid user ID: %v", err))
	}

	servers, err := DB.GetServers(uint(ID))
	if err != nil {
		log.Printf("Failed to get servers for user ID %d: %v", ID, err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get servers: %v", err))
	}

	if servers == nil {
		return &panel.ListServersResponse{}, nil
	}

	var serverDetails []*panel.Server
	for _, server := range servers {
		serverDetails = append(serverDetails, &panel.Server{
			Id:      server.ID,
			Name:    server.Name,
			Address: server.Address,
			Port:    server.Port,
			Status:  server.Status,
		})
	}

	return &panel.ListServersResponse{Servers: serverDetails}, nil
}

func (c *controlPanel) CreateDatabaseServer(ctx context.Context, req *panel.CreateDatabaseRequest) (*panel.CreateDatabaseResponse, error) {
	dbname, username, password := models.GenerateDatabaseCredentials()
	userID, _ := strconv.Atoi(req.UserId)
	err := database.CreateMySQLDatabase(DB.DBMySQL, uint(userID), req.Id, dbname, username, password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to create database server: %v", err))
	}

	return &panel.CreateDatabaseResponse{DatabaseId: dbname}, nil
}

func (c *controlPanel) GetDatabases(ctx context.Context, req *panel.GetDatabasesRequest) (*panel.GetDatabasesResponse, error) {
	userID, _ := strconv.Atoi(req.UserId)
	databases, err := database.GetMySqlDatabases(DB.DBMySQL, uint(userID), req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to get databases: %v", err))
	}

	var databaseDetails []*panel.Database
	for _, server := range databases {
		databaseDetails = append(databaseDetails, &panel.Database{
			DbName:   server.DBName,
			Address:  "localhost",
			Port:     "3036",
			Username: server.Username,
			Password: server.Password,
		})
	}

	return &panel.GetDatabasesResponse{Databases: databaseDetails}, nil
}

func (c *controlPanel) DeleteDatabase(ctx context.Context, req *panel.DeleteDatabaseRequest) (*panel.DeleteDatabaseResponse, error) {
	userID, _ := strconv.Atoi(req.UserId)
	dbname, err := database.DeleteMySqlDatabase(DB.DBMySQL, uint(userID), req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("failed to delete database: %v", err))
	}

	return &panel.DeleteDatabaseResponse{DatabaseId: dbname}, nil
}

func (c *controlPanel) UploadFiles(ctx context.Context, stream *panel.PanelService_UploadArchiveServer) error {

	return nil
}
