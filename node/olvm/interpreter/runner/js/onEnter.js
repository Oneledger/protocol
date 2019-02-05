function should(condition, errorObj) {
  if (condition) {
    return;
  }
  if(typeof errorObj == 'undefined') {
    errorObj = new Error('General Error');
  }
  throw errorObj;
}
var SafeMath = {
  NUMBER_IS_TOO_BIG: 1,
  NUMBER_IS_TOO_SMALL: 2,
  max: function() {
    return 1.157920892373162e+77;//2 ** 256
  },
  min: function() {
    return -1.157920892373162e+77;//- 2 ** 256
  },
  add: function(a, b) {
    var c = a + b
    should(c <= this.max(), new Error(this.NUMBER_IS_TOO_BIG));
    should(c >= this.min(),new Error(this.NUMBER_IS_TOO_SMALL));
    return c;
  },
  sub: function(a, b) {
    var c = a - b;
    should(c <= this.max(), new Error(this.NUMBER_IS_TOO_BIG));
    should(c >= this.min(),new Error(this.NUMBER_IS_TOO_SMALL));
    return c;
  },
  mul: function (a, b) {
    var c = a * b;
    should(c <= this.max(), new Error(this.NUMBER_IS_TOO_BIG));
    should(c >= this.min(),new Error(this.NUMBER_IS_TOO_SMALL));
    return c;
  },
  div: function (a, b) {
    var c = a / b;
    should(c <= this.max(), new Error(this.NUMBER_IS_TOO_BIG));
    should(c >= this.min(),new Error(this.NUMBER_IS_TOO_SMALL));
    return c;
  }
}

var IndexSet = function () { ///this is a simple porting for the set which lacks from ES5
  this._innerArray = [];
}

IndexSet.prototype.add = function(val) {
  for (var i = 0; i < this._innerArray.length; i ++) {
    if (this._innerArray[i] === val) {
      return;
    }
  }
  this._innerArray.push(val);
}

IndexSet.prototype.remove = function(val) {
  var replacedArray = [];
  for (var i = 0; i < this._innerArray.length; i ++) {
    if (this._innerArray[i] != val) {
      replacedArray.push(this._innerArray[i]);
    }
  }
  this._innerArray = replacedArray;
}

IndexSet.prototype.has = function (val) {
  for (var i = 0; i < this._innerArray.length; i ++) {
    if (this._innerArray[i] === val) {
      return true;
    }
  }
  return false;
}

IndexSet.prototype.list = function () {
  return this._innerArray;
}

var Context = function () {
  this.storage = {};
  this.line_data = [];//line data will be used to add the command execution steps for blockchain
  this.updateIndexSet = new IndexSet();
}

Context.prototype.get = function (key) {
  var val = this.storage[key];
  if (typeof val === 'undefined') {
    //try get val from context object
    val = __GetContextValue__(key);
    if (typeof val === undefined || val == null){
      return undefined;
    }
    val = JSON.parse(val);
  }
  return val;
}

Context.prototype.set = function (key, val) {
  this.updateIndexSet.add(key);
  this.storage[key] = val;
}

Context.prototype.getStorage = function () {
  return this.storage;
}

Context.prototype.getUpdateIndexList = function () {
  return this.updateIndexSet.list();
}
Context.prototype.execute =  function(action, parameters) {
  this.line_data.push({action: action, parameters: parameters});
}
Context.prototype.getLineData = function () {
  return this.line_data;
}

if (__callString__ == '') {
  __callString__ = "default__()";
}

var context = new Context();
