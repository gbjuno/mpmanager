/**
 * Created by Jingle on 2017/10/28.
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

/**
 * 传入一个Date对象转成可以传入URL的String对象
 * 例如： 2017/7/3 05:20:00  -> 20170703
 * @param {*} date 
 */
export const getDateQueryString = (date) => {
    let year = date.getFullYear()
    let month = (date.getMonth()+1).toString(); month = month.length === 1 ? '0' + month : month
    let day = date.getDate().toString(); day = day.length === 1 ? '0' + day : day
    return `${year}${month}${day}`
}