/**
 * Created by Jingle on 2017/12/10.
 */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as _ from 'lodash'
import moment from 'moment';
import * as download from 'downloadjs'
import * as config from '../../axios/config';
import { Form, Icon, Input, Button, Select, DatePicker, message } from 'antd';
import { fetchData, receiveData } from '../../action';
import MyDatePicker from '../../components/date-picker'

const FormItem = Form.Item;
const Search = Input.Search;
const Option = Select.Option;
const { YearPicker, MonthPicker } = DatePicker;

const dateFormat = 'YYYY-MM-DD';
const queryDateFormat = 'YYYYMMDD';

function hasErrors(fieldsError) {
    return Object.keys(fieldsError).some(field => fieldsError[field]);
}

class PhotoStatusSearch extends Component {

    state = {
        townsData: [],
        countriesData: [],
        selectedTownId: '',
        selectedCountryId: '',
        selectedDate: new Date(),
    }

    componentDidMount() {
        // To disabled submit button at the beginning.
        this.props.form.validateFields();
    }


    handleSubmit = (e) => {
        e.preventDefault();
        
        this.props.form.validateFields((err, values) => {
            if (!err) {
                const { fetchData } = this.props
                if(_.toString(values.town) === '0'){
                    fetchData({
                        funcName: 'fetchCompanies', params: {
                        }, stateName: 'companiesData'
                    })
                }
                if(_.isEmpty(values.country)){
                    message.info('请选择村')
                    return
                }
                fetchData({funcName: 'fetchCompaniesByCountryId', params: { 
                    countryId: values.country}, 
                    stateName: 'companiesData'})
            }
        });
    };


    downloadFile = () => {

        const { year, month } = this.state

        if(this.props.exportReport){
            this.props.exportReport(year, month)
        }
    }

    getOptions = ( ) => {

        const data=[
            {key: '01', id: '01', name: '一月'},
            {key: '02', id: '02', name: '二月'},
            {key: '03', id: '03', name: '三月'},
            {key: '04', id: '04', name: '四月'},
            {key: '05', id: '05', name: '五月'},
            {key: '06', id: '06', name: '六月'},
            {key: '07', id: '07', name: '七月'},
            {key: '08', id: '08', name: '八月'},
            {key: '09', id: '09', name: '九月'},
            {key: '10', id: '10', name: '十月'},
            {key: '11', id: '11', name: '十一月'},
            {key: '12', id: '12', name: '十二月'},
        ] 
        return data.map(item => {
            return <Option key={item.key} value={`${item.id}`}>{item.name}</Option>
        })
    }

    onYearChange = (value, dateStr) => {
        if(_.isEmpty(value )){
            this.setState({
                year: undefined,
            })
            return
        }
        this.setState({
            year: value.format('YYYY')
        })
    }

    onMonthChange = ( value ) => {
        this.setState({
            month: value,
        })
    }


    render() {
        const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched } = this.props.form;
        const { style, filter } = this.props
        const { townsData, countriesData, year } = this.state


        // Only show error after a field is touched.
        const fileNameError = isFieldTouched('fileName') && getFieldError('fileName');
        return (
            <Form layout="inline" style={style} onSubmit={this.handleSubmit}>
                <FormItem
                    style={{paddingBottom: 13}}
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    <MyDatePicker  topMode="year" format="YYYY" onChange={this.onYearChange}/>
                </FormItem>
                <FormItem
                    style={{paddingBottom: 13}}
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                        <Select
                        showSearch
                        style={{ width: 200 }}
                        placeholder="请选择月份"
                        optionFilterProp="children"
                        onChange={this.onMonthChange}
                        filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                        >
                            {this.getOptions()}
                        </Select>
                </FormItem>
                <FormItem 
                    style={{paddingBottom: 13}}
                >
                    <Button
                        onClick={this.downloadFile}
                        type="primary"
                    >
                       下载报表
                    </Button>
                </FormItem>
            </Form>
        );
    }
}

const mapStateToProps = state => {
    const { searchFilter } = state
    return { ...state.httpData, filter: searchFilter};
};
const mapDispatchToProps = dispatch => ({
    receiveData: bindActionCreators(receiveData, dispatch),
    fetchData: bindActionCreators(fetchData, dispatch),
});

export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(PhotoStatusSearch))
