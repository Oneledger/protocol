class Logger {
    info(any, ...msg){
        console.log(new Date().toISOString() + ' | ' + any, ...msg);
    }
    warn(any, ...msg){
        console.warn(new Date().toISOString() + ' | ' + anyany, ...msg);
    }
    error(any, ...msg){
        console.error(new Date().toISOString() + ' | ' + any, ...msg);
    }
    devInfo(...args) {
        if (process.env.NODE_ENV === 'development') {
            this.info(...args);
        }
    }
}

exports.logger = new Logger();
