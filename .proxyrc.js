const allowedOrigins = ['http://localhost:3001', 'http://localhost:1234']

module.exports = function (app) {
    app.use((req, res, next) => {
        const origin = req.headers.origin;
        if (allowedOrigins.includes(origin)) {
            res.setHeader('Access-Control-Allow-Origin', origin);
        }

        res.setHeader('Access-Control-Allow-Methods','GET,POST,PUT,DELETE');
        res.setHeader('Access-Control-Allow-Headers','Content-Type','Authorization');
        res.setHeader('Access-Control-Allow-Credentials', true);

        next();
    });
};