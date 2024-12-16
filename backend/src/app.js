require("dotenv").config();
const express = require("express");
const { Server } = require("socket.io");
const http = require("http");

const app = express();
const server = http.createServer(app);
const io = new Server(server);

// Middleware
app.use(express.json());

// Basic Route
app.get("/", (req, res) => {
  res.send("AggieCode Backend is running!");
});

// WebSocket Setup
io.on("connection", (socket) => {
  console.log("A user connected");
  socket.on("disconnect", () => {
    console.log("A user disconnected");
  });
});

// Start the Server
const PORT = process.env.PORT || 3000;
server.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});

// Testing Database Connection
const db = require("./models/db");

app.get("/test-db", async (req, res) => {
  try {
    const result = await db.query("SELECT NOW()");
    res.json(result.rows[0]);
  } catch (err) {
    console.error(err);
    res.status(500).send("Database connection failed");
  }
});

// Testing Redis Connection
const redisClient = require("./models/redis");

app.get("/test-redis", async (req, res) => {
  try {
    await redisClient.set("test", "Hello from Redis!");
    const value = await redisClient.get("test");
    res.send(value);
  } catch (err) {
    console.error(err);
    res.status(500).send("Redis connection failed");
  }
});

