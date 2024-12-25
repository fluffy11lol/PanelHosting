const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const express = require('express');
const bodyParser = require('body-parser');

// Пути к .proto файлам
const ROOT_PROTO_PATH = './proto';
const PANEL_PROTO_PATH = `${ROOT_PROTO_PATH}/panel/panel.proto`;

// Загрузка .proto
const packageDefinition = protoLoader.loadSync(PANEL_PROTO_PATH, {
  includeDirs: [ROOT_PROTO_PATH],
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});
const protoDescriptor = grpc.loadPackageDefinition(packageDefinition);

// Доступ к описанию пакета panel
const panel = protoDescriptor.panel.v1;
console.log(panel);

// Создание клиента
const client = new panel.PanelService('localhost:50051', grpc.credentials.createInsecure());

// Настройка Express
const app = express();
app.use(bodyParser.json());

// HTTP-эндпоинт для вызова метода Login
app.post('/v1/login', (req, res) => {
  const { username, password } = req.body;

  // gRPC-вызов
  client.Login({ username, password }, (error, response) => {
    if (error) {
      console.error('gRPC Error:', error);
      return res.status(500).json({ error: error.message });
    }

    console.log('gRPC Response:', response);
    return res.json(response); // Отправляем ответ на фронтенд
  });
});

// Запуск сервера
const PORT = 3000;
app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
