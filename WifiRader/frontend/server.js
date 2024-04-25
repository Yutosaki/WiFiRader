require('dotenv').config();
const express = require('express');
const app = express();
const port = 3000;
const googleMapsApiKey = process.env.GOOGLE_MAPS_API_KEY;

const cors = require('cors');
app.use(express.static('../frontend'));
app.use(cors());

app.get('/', (req, res) => {
    res.send(`
    <!DOCTYPE html>
    <html lang="ja">
    <head>
        <meta charset="UTF-8">
        <title>Wifi Radar</title>
        <link rel="stylesheet" href="/app.css">
        <script async defer src="https://maps.googleapis.com/maps/api/js?key=${googleMapsApiKey}&language=ja&libraries=geometry"></script>
        <script src="/app.js" defer></script>
    </head>
    <body>
        <div id="container">
            <div id="sidebar">
                <button id="toggleButton">⇄</button>
                <input type="number" id="desiredAmount" placeholder="希望金額を入力" />
                <button onclick="submitLocationAndAmount()">送信</button>
                <iframe id="placeIframe"></iframe>
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