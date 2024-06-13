const util = require("node:util");
exports.PLATFORM = require("node:os").platform();

const STRING_LIST_SPLITTER = /\s*[\n\r\,;]+\s*/;

exports.stringlist = function (str = [], opts = {splitter: STRING_LIST_SPLITTER}) {
	if (typeof(str) === 'string') str = str.split(opts.splitter);
	// if (!Array.isArray(str)) str = [str]
	return str.map(line => line.trim()).filter(line => line !== '');
}

exports.inspect = (o, opts = {}) => util.inspect(o, {
	depth: 3,
	colors: true,
	compact: true,
	// numericSeparator: true,
	...opts,
})