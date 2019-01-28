import numeral from 'numeral';
import './g2';
import Bar from './Bar';

const yuan = val => `¥ ${numeral(val).format('0,0')}`;

const Charts = {
  yuan,
  Bar,
};

export {
  Charts as default,
  yuan,
  Bar,
};
