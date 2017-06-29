

var app = app || {};

(function (app) {

    var bus = {
        init: function (URL) {
            this._ws = new WebSocket(URL);
            this._ws.onopen = function () { console.log('opened'); }
            this._ws.onclose = function () { console.log('closed'); }
            this._ws.onmessage = function (data) {
                console.log("data", data)
            };
        },
        send: function (data) {
            this._ws.send(data);
        }
    }

    app.bus = bus;

})(app);