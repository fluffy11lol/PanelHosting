const protoLoader = require('@grpc/proto-loader');
const grpc = require('@grpc/grpc-js');

const ROOT_PROTO_PATH = './proto';

const packageDefinition = protoLoader.loadSync(
  [`${ROOT_PROTO_PATH}/panel/panel.proto`],
  {
    includeDirs: [ROOT_PROTO_PATH],
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
  }
);

const panelProto = grpc.loadPackageDefinition(packageDefinition).panel.v1;


function login(call, callback) {
    callback(null, {username: call.request.username, password: call.request.password});
}

function listServers(call, callback) {
  const servers = [
    { name: 'Server 1', address: '192.168.0.1', port: '8080', id: '1', status: 'online' },
    { name: 'Server 2', address: '192.168.0.2', port: '8081', id: '2', status: 'offline' }
  ];
  
  callback(null, { servers });
}

function getServerDetails(call, callback) {
  const server = { name: 'Server 1', address: '192.168.0.1', port: '8080', id: '1', status: 'online' };
  
  if (server.id === call.request.id) {
    callback(null, { server });
  } else {
    callback({
      code: grpc.status.NOT_FOUND,
      details: 'Server not found'
    });
  }
}

function createServer(call, callback) {
  const newServer = {
    name: call.request.name,
    address: '192.168.0.3',
    port: '8083',
    id: '3',
    status: 'offline'
  };
  
  callback(null, { server: newServer });
}

function updateServer(call, callback) {
  const updatedServer = {
    name: call.request.name,
    address: '192.168.0.1',
    port: '8080',
    id: call.request.id,
    status: 'online'
  };

  callback(null, { server: updatedServer });
}

function deleteServer(call, callback) {
  callback(null, { id: call.request.id });
}

function startServer() {
  const server = new grpc.Server();

  server.addService(panelProto.PanelService.service, {
    Login: login,
    ListServers: listServers,
    GetServerDetails: getServerDetails,
    CreateServer: createServer,
    UpdateServer: updateServer,
    DeleteServer: deleteServer
  });

  server.bindAsync('0.0.0.0:50051', grpc.ServerCredentials.createInsecure(), () => {
    console.log('Server running at http://127.0.0.1:50051');
    server.start();
  });
}

startServer();
