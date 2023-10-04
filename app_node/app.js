const express = require('express');
const bodyParser = require('body-parser');

const app = express();
const port = 3000;

app.use(bodyParser.json());

app.get('/', (req, res) => {
  const response = {
    response: 'This is from Node.js Backend, Hii Pintu',
  };
  res.json(response);
});

app.listen(port, () => {
  console.log(`[bootup]: Server is running at port: ${port}`);
});
