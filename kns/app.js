require('dotenv').config();

const express = require('express');

const app = express();

const searchRoutes = require('./routes/search');

// routes
app.use('/api/v1/search', searchRoutes);


// listening the server
app.listen(process.env.PORT, err => {
    if (err) console.error(err);

    console.log('Server listening to port :', process.env.PORT);
});