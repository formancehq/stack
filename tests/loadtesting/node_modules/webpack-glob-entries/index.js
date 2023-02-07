var glob = require('glob'),
    path = require('path');

module.exports = function (globPath) {
    var files = glob.sync(globPath);
    var entries = {};

    for (var i = 0; i < files.length; i++) {
        var entry = files[i];
        entries[path.basename(entry, path.extname(entry))] = entry;
    }

    return entries;
}
