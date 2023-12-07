const value = process.argv[2];

const getValue = (string) => {
    const value = string.match(/([a-z].*)\/([a-z].*)\/([a-z].*)/)[2];
    return value;
};

console.log(getValue(value));
