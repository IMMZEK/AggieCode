// Note: This might be replaced with Firebase Auth in the future.
import axios from 'axios';

// Constants for OAuth2 configuration
const GOOGLE_OAUTH2_ENDPOINT = 'https://accounts.google.com/o/oauth2/v2/auth';
const CLIENT_ID = '544926959381-m1t3pej8jk0qpohdu6du2i7t5r71jhdh.apps.googleusercontent.com';
const REDIRECT_URI = 'http://localhost:5173/oauth2callback';
const SCOPES = ['https://www.googleapis.com/auth/devstorage.read_write'];

// TODO: Implement / change the URls below
const BACKEND_TOKEN_ENDPOINT = 'http://localhost:3000/api/token';
const BACKEND_REFRESH_ENDPOINT = 'http://localhost:3000/api/refresh';

const OAUTH_STATE_KEY = 'oauth2_state';

export default {
    // Generates the Google OAuth2 authorization URL
    generateAuthUrl() {
        const state = Math.random().toString(36).substring(2, 15); // Generate a random state string
        localStorage.setItem(OAUTH_STATE_KEY, state); // Store the state in localStorage

        const url = new URL(GOOGLE_OAUTH2_ENDPOINT);
        url.searchParams.set('client_id', CLIENT_ID);
        url.searchParams.set('redirect_uri', REDIRECT_URI);
        url.searchParams.set('response_type', 'code');
        url.searchParams.set('scope', SCOPES.join(' '));
        url.searchParams.set('state', state);
        url.searchParams.set('access_type', 'offline'); // Important for refresh token

        return url.toString(); // Return the complete URL as a string
    },

    // Exchanges the authorization code for an access token
    async exchangeCodeForToken(code, state) {
        const storedState = localStorage.getItem(OAUTH_STATE_KEY); // Retrieve the stored state
        localStorage.removeItem(OAUTH_STATE_KEY); // Remove the state from localStorage

        if (state !== storedState) {
            throw new Error('State mismatch error'); // Throw an error if the state doesn't match
        }

        try {
            // Call your backend's /api/token endpoint
            const response = await axios.post(BACKEND_TOKEN_ENDPOINT, {
                code,
                redirect_uri: REDIRECT_URI,
            });

            return response.data; // Return the token data { access_token, refresh_token, expires_in }
        } catch (error) {
            console.error(
                'Error exchanging code for token:',
                error.response.data || error.message
            );
            throw error; // Rethrow the error
        }
    },

    // Refreshes the access token using the refresh token
    async refreshAccessToken(refreshToken) {
        try {
            // Call your backend's /api/refresh endpoint
            const response = await axios.post(BACKEND_REFRESH_ENDPOINT, {
                refreshToken,
            });

            return response.data; // Return the new token data { access_token, expires_in }
        } catch (error) {
            console.error(
                'Error refreshing access token:',
                error.response.data || error.message
            );
            throw error; // Rethrow the error
        }
    },
};