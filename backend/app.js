require('dotenv').config();
const express = require('express');
const axios = require('axios');
const cors = require('cors');
const app = express();

// CORS - Allowed frontend origin
const allowedOrigins = [process.env.FRONTEND_ORIGIN];
const corsOptions = {
    origin: function (origin, callback) {
        if (allowedOrigins.indexOf(origin) !== -1 || !origin) {
            callback(null, true);
        } else {
            callback(new Error('Not allowed by CORS'));
        }
    },
};

// Midldlewares
app.use(cors(corsOptions));
app.use(express.json());

const GOOGLE_TOKEN_ENDPOINT = 'https://oauth2.googleapis.com/token';
const CLIENT_ID = process.env.GOOGLE_CLIENT_ID;
const CLIENT_SECRET = process.env.GOOGLE_CLIENT_SECRET;

// Endpoint to exchange code for token
app.post('/api/token', async (req, res) => {
    const { code, redirect_uri } = req.body;

    // Exchange code for token
    try {
        const response = await axios.post(GOOGLE_TOKEN_ENDPOINT, {
            code,
            client_id: CLIENT_ID,
            client_secret: CLIENT_SECRET,
            redirect_uri,
            grant_type: 'authorization_code',
        });

        res.json(response.data);
    } catch (error) {
        console.error('Error exchanging code for token:', error.response.data || error.message);
        res.status(500).json({ error: 'Failed to exchange code for token' });
    }
});

// Endpoint to refresh token
app.post('/api/refresh', async (req, res) => {
    const { refreshToken } = req.body;

    // Refresh token
    try {
        const response = await axios.post(GOOGLE_TOKEN_ENDPOINT, {
            refresh_token: refreshToken,
            client_id: CLIENT_ID,
            client_secret: CLIENT_SECRET,
            grant_type: 'refresh_token',
        });

        res.json(response.data);
    } catch (error) {
        console.error('Error refreshing token:', error.response.data || error.message);
        res.status(500).json({ error: 'Failed to refresh token' });
    }
});

// Start server
const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    console.log(`Backend server running on port ${PORT}`);
});