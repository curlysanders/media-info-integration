(function (global) {
    "use strict";

    var namespace = global.SYNO = global.SYNO || {};
    namespace.MediaInfoPoc = namespace.MediaInfoPoc || {};

    var supportedExtensions = {
        "m4v": true,
        "mkv": true,
        "mov": true,
        "mp4": true,
        "mpeg": true,
        "mpg": true,
        "webm": true
    };

    function firstRecord(records) {
        if (Array.isArray(records)) {
            return records.length === 1 ? records[0] : null;
        }
        if (records && typeof records.getCount === "function" && typeof records.getAt === "function") {
            return records.getCount() === 1 ? records.getAt(0) : null;
        }
        if (records && Array.isArray(records.records)) {
            return records.records.length === 1 ? records.records[0] : null;
        }
        if (records && Array.isArray(records.recs)) {
            return records.recs.length === 1 ? records.recs[0] : null;
        }
        if (records && typeof records === "object") {
            return records;
        }
        return null;
    }

    function value(record, name) {
        if (!record) {
            return null;
        }
        if (typeof record.get === "function") {
            return record.get(name);
        }
        if (record.data && Object.prototype.hasOwnProperty.call(record.data, name)) {
            return record.data[name];
        }
        return record[name];
    }

    function selection(records) {
        var record = firstRecord(records);
        var path = value(record, "path");
        var name = value(record, "name");
        if (!record || value(record, "isdir") || typeof path !== "string" || typeof name !== "string") {
            return null;
        }
        var extension = name.split(".").pop().toLowerCase();
        if (!Object.prototype.hasOwnProperty.call(supportedExtensions, extension)) {
            return null;
        }
        return { name: name, path: path };
    }

    namespace.MediaInfoPoc.FileStation = {
        isEligible: function (records) {
            return selection(records) !== null;
        },

        launch: function (records) {
            var selected = selection(records);
            if (!selected) {
                return false;
            }

            var popup = global.open(
                "/webman/3rdparty/mediainfo-poc/index.html",
                "mediainfo-poc",
                "width=760,height=520"
            );
            if (!popup) {
                return false;
            }

            var sendSelection = function (event) {
                if (event.source !== popup || event.origin !== global.location.origin || !event.data || event.data.type !== "mediainfo-poc-ready") {
                    return;
                }
                popup.postMessage({
                    type: "mediainfo-poc-selection",
                    selection: selected
                }, global.location.origin);
                global.removeEventListener("message", sendSelection);
            };
            global.addEventListener("message", sendSelection);
            global.setTimeout(function () {
                global.removeEventListener("message", sendSelection);
            }, 5000);
            return true;
        }
    };
}(window));
