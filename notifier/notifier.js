const express = require("express");
const { notify } = require("./notify");

const app = express();

app.use(express.json());

app.get("/health", (req, res) => res.status(200).send());

app.post("/notify", (req, res) => {
  notify(req.body, (reply) => res.send(reply));
});

const port = process.env.PORT || 5000;

app.listen(port, () => {
  console.log(`Server is runnig on port: ${port}`);
});
