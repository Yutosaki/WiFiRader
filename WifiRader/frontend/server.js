const express = require('express');
const app = express();
const port = 3000;

// 環境変数から API キーを取得
const googleMapsApiKey = process.env.GOOGLE_MAPS_API_KEY;

app.use(express.static('public'));

app.get('/', (req, res) => {
    // HTMLファイルに API キーを埋め込んでレンダリング
    res.send(`<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <title>Wifi Rader</title>
    <link rel="stylesheet" href="app.css">
    <script src="https://maps.googleapis.com/maps/api/js?key=${googleMapsApiKey}&language=ja"></script>
    <script src="app.js" defer></script>
</head>
<body>
    <div id="map"></div>
</body>
</html>`);
});

app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
});