/**
 * Created by hao.cheng on 2017/4/28.
 */
// 获取url的参数
export const queryString = () => {
    let _queryString = {};
    const _query = window.location.search.substr(1);
    const _vars = _query.split('&');
    _vars.forEach((v, i) => {
        const _pair = v.split('=');
        if (!_queryString.hasOwnProperty(_pair[0])) {
            _queryString[_pair[0]] = decodeURIComponent(_pair[1]);
        } else if (typeof _queryString[_pair[0]] === 'string') {
            const _arr = [ _queryString[_pair[0]], decodeURIComponent(_pair[1])];
            _queryString[_pair[0]] = _arr;
        } else {
            _queryString[_pair[0]].push(decodeURIComponent(_pair[1]));
        }
    });
    return _queryString;
};

export const number =  (data, n) =>{
    let numbers = '';
    for (var i = 0; i < n; i++) {
        numbers += '0';
    }
    let s = 1 + numbers;
    // 如果是整数需要添加后面的0
    let spot = "." + numbers;
    // Math.round四舍五入  
    //  parseFloat() 函数可解析一个字符串，并返回一个浮点数。
    let value = Math.round(parseFloat(data) * s) / s;
    // 从小数点后面进行分割
    let d = value.toString().split(".");
    if (d.length == 1) {
        value = value.toString() + spot;
        return value;
    }
    if (d.length > 1) {
        if (d[1].length < 2) {
            value = value.toString() + "0";
        }
        return value;
    }
}
