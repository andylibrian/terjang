process.env.VUE_APP_SERVER_BASE_URL = process.env.VUE_APP_SERVER_BASE_URL ?? process.env.NODE_ENV === 'production' ? '' : 'http://localhost:9009';

module.exports = {
    pages: {
        index: {
            entry: './src/main.js',
            title: "Terjang"
        }
    }
}
