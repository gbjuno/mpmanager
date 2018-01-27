/**
 * Created by Jingle on 2017/12/10.
 */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as download from 'downloadjs';
import * as _ from 'lodash'
import moment from 'moment';
import { Form, Icon, Input, Button, Select, DatePicker, Row, message } from 'antd';
import * as CONSTANTS from '../../constants';
import { fetchData, receiveData } from '../../action';
import * as config from '../../axios/config';

const FormItem = Form.Item;
const Search = Input.Search;
const Option = Select.Option;

function hasErrors(fieldsError) {
    return Object.keys(fieldsError).some(field => fieldsError[field]);
}

class SummarySearch extends Component {

    state = {
        townsData: [],
        countriesData: [],
        companiesData: [],
        selectedTownId: '',
        selectedCountryId: '',
        selectedCompanyId: '',
        currentDate: new Date(),
        fromDate: null,
        toDate: null,
        isFirst: true,  //未点击搜索按钮
        defaultCompanyId: '',
    }

    componentDidMount() {
        // To disabled submit button at the beginning.
        this.props.form.validateFields();
        this.fetchTownList();
        let from = moment(this.state.currentDate).format(CONSTANTS.DATE_QUERY_FORMAT)
        let to = from
        this.searchSummary('', from, to)
    }


    handleSubmit = (e) => {
        e.preventDefault();
        const { isFirst, defaultCompanyId } = this.state
        this.props.form.validateFields((err, values) => {
            if (!err) {
                if(values.fromDate === null && values.toDate === null) {
                    message.error('请选择开始日期或结束日期')
                    return
                }
                let fromDate = values.fromDate.format(CONSTANTS.DATE_QUERY_FORMAT);
                let toDate = values.toDate.format(CONSTANTS.DATE_QUERY_FORMAT);

                if(fromDate > toDate){
                    message.error('开始日期不能大于结束日期')
                    return
                }

                let companyId = values.company
                if(companyId === undefined)companyId=''
                this.searchSummary(companyId, fromDate, toDate)
            }
        });
    };

    searchSummary = (companyId, fromDate, toDate) => {
        const { fetchData } = this.props
        fetchData({funcName: 'searchSummaries', params: {companyId, from: fromDate, to: toDate}, 
            stateName: 'summariesData'})
    }

    onFromDateChange = (date, dateString) => {
        const { searchPicture } = this.props
        this.setState({
            isFirst: false,
            fromDate: date,
        })
        if (date === undefined || date === null) return
        
    }

    onToDateChange = (date, dateString) => {
        const { searchPicture } = this.props
        this.setState({
            isFirst: false,
            toDate: date,
        })
        if (date === undefined || date === null) return
        
    }

    onTownChange = (value) => {
        const { form } = this.props
        this.setState({
            selectedTownId: value,
            isFirst: false,
        },() => this.fetchCountryList(value))
        form.setFieldsValue({
            country: undefined,
            company: undefined,
        })
    }


    onCountryChange = (value) => {
        const { form } = this.props
        this.setState({
            selectedCountryId: value,
            isFirst: false,
        },() => this.fetchCompanyList(value))
        form.setFieldsValue({
            company: undefined,
        })
    }

    onCompanyChange = (value) => {
        this.setState({
            isFirst: false,
            selectedCompanyId: value,
        })
    }


    fetchTownList = () => {
        const { fetchData } = this.props
        fetchData({funcName: 'fetchTowns', stateName: 'townsData'}).then(res => {
            if(res === undefined || res.data === undefined || res.data.towns === undefined) return
            this.setState({
                townsData: [...res.data.towns.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            }, () => {
                // this.fetchCountryList(res.data.towns[0].id)
            });
        });
    }

    fetchCountryList = (selectedTownId) => {
        const { fetchData } = this.props
        if(selectedTownId === undefined){
            return
        }
        fetchData({funcName: 'fetchCountries', stateName: 'countriesData', params: {townId: selectedTownId}}).then(res => {
            if(res === undefined || res.data === undefined || res.data.countries === undefined) return
            this.setState({
                countriesData: [...res.data.countries.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            }, () => {
                // this.fetchCompanyList(res.data.countries[0].id)
            });
        });
    }

    fetchCompanyList = (selectedCountryId) => {
        const { fetchData } = this.props
        const { isFirst } = this.state
        if(selectedCountryId === undefined){
            return
        }
        fetchData({funcName: 'fetchCompaniesByCountryId', stateName: 'companiesData', params: {countryId: selectedCountryId}}).then(res => {
            if(res === undefined || res.data === undefined || res.data.companies === undefined) return
            this.setState({
                companiesData: [...res.data.companies.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            }, () => {
                if( isFirst ){
                    
                }
            });
        });
    }


    getOptions = ( data=[] ) => {
        
        return data.map(item => {
            return <Option key={item.key} value={item.id}>{item.name}</Option>
        })
    }

    genDay = (date) =>{
        let defaultDay = moment(this.state.currentDate).format(CONSTANTS.DATE_QUERY_FORMAT)
        if(date === null){
            return defaultDay
        }else{
            return date.format(CONSTANTS.DATE_QUERY_FORMAT)
        }
    }

    downloadReport = () => {
        const { selectedCompanyId, fromDate, toDate } = this.state
        let from, to;
        from = this.genDay(fromDate)
        to = this.genDay(toDate)

        if(from > to){
            message.error('开始日期不能大于结束日期')
            return
        }

        const filter = {
            companyId: selectedCompanyId,
            from,
            to,
        }
        let url = config.EXPORT_SUMMARY_URL(filter); //TODO: 换成下载公司数据url,及相应的文件格式
        const x = new XMLHttpRequest;
        x.open("GET", url, true);
        x.responseType = "blob";
        x.onload = function (e) {
            download(x.response, "统计报表.xlsx", "application/octet-stream")
        }
        x.send();
    }


    render() {
        const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched } = this.props.form;
        const { style, filter } = this.props
        const { townsData, countriesData, companiesData, currentDate } = this.state
        let fromDate = currentDate, toDate = currentDate;

        // Only show error after a field is touched.
        const fileNameError = isFieldTouched('fileName') && getFieldError('fileName');
        return (
            <Form layout="inline" style={style} onSubmit={this.handleSubmit}>
                <Row gutter={24}>
                <FormItem 
                    style={{paddingBottom: 13}}
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('town', {
                        // initialValue: townsData[0]? townsData[0].id:'',
                    })(
                        <Select
                        showSearch
                        style={{ width: 200 }}
                        placeholder="请选择镇"
                        optionFilterProp="children"
                        onChange={this.onTownChange}
                        filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                        >
                            {this.getOptions(townsData)}
                        </Select>
                    )}
                </FormItem>
                <FormItem
                    style={{paddingBottom: 13}}
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('country', {
                        // initialValue: countriesData[0]? countriesData[0].id:'',
                    })(
                        <Select
                        showSearch
                        style={{ width: 200 }}
                        placeholder="请选择村"
                        optionFilterProp="children"
                        onChange={this.onCountryChange}
                        filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                        >
                            {this.getOptions(countriesData)}
                        </Select>
                    )}
                </FormItem>
                <FormItem
                    style={{paddingBottom: 13}}
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('company', {
                        // initialValue: companiesData[0]? companiesData[0].id:'',
                        rule: [
                            {require: true},
                        ]
                    })(
                        <Select
                        showSearch
                        style={{ width: 200 }}
                        placeholder="请选择公司"
                        optionFilterProp="children"
                        onChange={this.onCompanyChange}
                        onSelect={this.onCompanySelect}
                        filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                        >
                            {this.getOptions(companiesData)}
                        </Select>
                    )}
                </FormItem>
                </Row>
                <Row gutter={24}>
                <FormItem
                    label="从"
                    style={{paddingBottom: 13}}
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('fromDate', {
                        initialValue: moment(fromDate, CONSTANTS.DATE_DISPLAY_FORMAT)
                    })(
                        <DatePicker onChange={this.onFromDateChange}/>
                    )}
                </FormItem>
                <FormItem
                    label="到"
                    style={{paddingBottom: 13}}
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('toDate', {
                        initialValue: moment(toDate, CONSTANTS.DATE_DISPLAY_FORMAT)
                    })(
                        <DatePicker onChange={this.onToDateChange}/>
                    )}
                </FormItem>
                <FormItem 
                    style={{paddingBottom: 13}}
                >
                    <Button
                        type="primary"
                        htmlType="submit"
                    >
                       搜索
                    </Button>
                    <Button onClick={this.downloadReport}
                        type="primary"
                    >
                        导出
                    </Button>
                </FormItem>
                </Row>
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

export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(SummarySearch))
