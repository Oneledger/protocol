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
  this.updateIndexSet = new IndexSet();
}

Context.prototype.get = function (key) {
  return this.storage[key];
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

module.exports = new Context();
