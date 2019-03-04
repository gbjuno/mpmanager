/**
 * Created by Jingle on 2017/12/10.
 */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as _ from 'lodash'
import moment from 'moment';
import { Form, Icon, Input, Button, Select, DatePicker, message } from 'antd';
import * as CONSTANTS from '../../constants';
import { fetchData, receiveData } from '../../action';

const FormItem = Form.Item;
const Search = Input.Search;
const Option = Select.Option;

function hasErrors(fieldsError) {
    return Object.keys(fieldsError).some(field => fieldsError[field]);
}

class VacationSearch extends Component {

    state = {
        townsData: [],
        countriesData: [],
        companiesData: [],
        selectedTownId: '',
        selectedCountryId: '',
        isFirst: true,  //未点击搜索按钮
    }

    componentDidMount() {
        // To disabled submit button at the beginning.
        this.props.form.validateFields();
        this.fetchCompanyList()
    }


    handleSubmit = (e) => {
        e.preventDefault();
        
        this.props.form.validateFields((err, values) => {
            if (!err) {
                let companyId = values.company
                if(companyId === undefined) {
                    message.error('请选择公司')
                    return
                }
                this.searchVacation(companyId)
            }
        });
    };

    searchVacation = (companyId) => {
        const { fetchData, onSearch } = this.props
        fetchData({funcName: 'searchVacations', params: {companyId}, 
            stateName: 'vacationsData'})
            .then(res => {
                if(onSearch){
                    onSearch()
                }
            })
    }



    onCompanySelect = (value, option) => {
        this.setState({
            isFirst: false,
        })
        if(this.props.onChange){
            let companyName = option.props.children
            this.props.onChange(value, companyName)
        }
    }




    fetchCompanyList = () => {
        const { fetchData } = this.props
        const { isFirst } = this.state

        fetchData({funcName: 'fetchCompanies', stateName: 'companiesData', params: {}}).then(res => {
            if(res === undefined || res.data === undefined || res.data.companies === undefined) return
            this.setState({
                companiesData: [...res.data.companies.map(val => {
                    val.key = val.id;
                    return val;
                })],
                loading: false,
            }, () => {
                if(isFirst){
                    // this.searchPlace(res.data.companies[0].id)
                    // if(this.props.onChange){
                    //     this.props.onChange(res.data.companies[0].id, res.data.companies[0].name)
                    // }
                }
            });
        });
    }


    getOptions = ( data=[] ) => {

        if(data && data.length > 0 && data[0].key !== 0){
            let startElement = {
                key: 0,
                id: 0,
                name: '全局',
            }
            data.unshift(startElement)
        }
        return data.map(item => {
            return <Option key={item.key} value={item.id} title={item.name}>{item.name}</Option>
        })
    }


    render() {
        const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched } = this.props.form;
        const { style, filter } = this.props
        const { townsData, countriesData, companiesData, selectedDate } = this.state

        // Only show error after a field is touched.
        const fileNameError = isFieldTouched('fileName') && getFieldError('fileName');
        return (
            <Form layout="inline" style={style} onSubmit={this.handleSubmit}>
                <FormItem
                    style={{paddingBottom: 13}}
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('company', {
                        initialValue: companiesData[0]? companiesData[0].id:'',
                        rule: [
                            {require: true},
                        ]
                    })(
                        <Select
                        showSearch
                        style={{ width: 300 }}
                        placeholder="请选择公司"
                        optionFilterProp="children"
                        onSelect={this.onCompanySelect}
                        filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                        >
                            {this.getOptions(companiesData)}
                        </Select>
                    )}
                </FormItem>
                {/* <FormItem 
                    style={{paddingBottom: 13}}
                >
                    <Button
                        type="primary"
                        htmlType="submit"
                    >
                       搜索
                    </Button>
                </FormItem> */}
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

export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(VacationSearch))
