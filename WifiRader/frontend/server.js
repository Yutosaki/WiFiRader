require('dotenv').config();
const express = require('express');
const app = express();
const port = 3000;
const googleMapsApiKey = process.env.GOOGLE_MAPS_API_KEY;

app.use(express.static('frontend'));

app.get('/', (req, res) => {
    res.send(`
    <!DOCTYPE html>
    <html lang="ja">
    <head>
        <meta charset="UTF-8">
        <title>Wifi Radar</title>
        <link rel="stylesheet" href="app.css">
        <script src="https://maps.googleapis.com/maps/api/js?key=${googleMapsApiKey}&language=ja"></script>
        <script src="app.js" defer></script>
    </head>
    <body>
        <div id="container">
            <div id="sidebar">
                console.log('Google Maps API Key:', process.env.GOOGLE_MAPS_API_KEY);
                <button id="toggleButton">⇄</button>
                <input type="number" id="desiredAmount" placeholder="希望金額を入力" />
                <button onclick="submitLocationAndAmount()">送信</button>
            </div>
            <div id="map"></div>
        </div>
    </body>
    </html>
    `);
});

const startServer = async () => {
    app.listen(port, async () => {
        console.log(`Server running on http://localhost:${port}`);
        try {
            const open = (await import('open')).default;
            open(`http://localhost:${port}`);
        } catch (error) {
            console.error("Failed to open browser:", error);
        }
    });
};

startServer();