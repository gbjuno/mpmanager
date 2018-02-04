/**
 * Created by Jingle on 2017/12/10.
 */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as _ from 'lodash'
import moment from 'moment';
import { Form, Icon, Input, Button, Select, DatePicker, message } from 'antd';
import { fetchData, receiveData } from '../../action';

const FormItem = Form.Item;
const Search = Input.Search;
const Option = Select.Option;

const dateFormat = 'YYYY-MM-DD';
const queryDateFormat = 'YYYYMMDD';

function hasErrors(fieldsError) {
    return Object.keys(fieldsError).some(field => fieldsError[field]);
}

const ACTION = {
    CREATE: 'create',
    UPDATE: 'update',
}

class PlaceForm extends Component {

    state = {
        value: this.props.value,
        action: this.props.action || ACTION.CREATE,
        townsData: [],
        countriesData: [],
        companiesData: [],
        selectedTownId: '',
        selectedCountryId: '',
        selectedDate: new Date(),
    }

    componentDidMount() {
        // To disabled submit button at the beginning.
        //this.props.form.validateFields();
    }


    handleSubmit = (e) => {
        e.preventDefault();
        const { value } = this.state
        this.props.form.validateFields((err, values) => {
            if (!err) {
                const { fetchData } = this.props
                let saveObj = {
                    name: values.name,
                    company_id: parseInt(values.company_id),
                    monitor_type_id: value.monitor_type_id,
                }
                
            }
        });
    };

    handleChange = (e) => {
        if(this.props.onChange){
            this.props.onChange(e.target.value)
        }
    }


    render() {
        const { getFieldDecorator } = this.props.form;
        const { style, filter } = this.props
        const { value } = this.state

        return (
            <Form style={style} onSubmit={this.handleSubmit}>
                <FormItem 
                    style={{}}
                >
                    {getFieldDecorator('companyName', {
                        rules: [{
                            required: true, message: '地点必填!',
                        }],
                    })(
                        <p>公司：{value.company_name}</p>
                    )}
                </FormItem>
                <FormItem 
                    style={{}}
                >
                    {getFieldDecorator('name', {
                        rules: [{
                            required: true, message: '地点必填!',
                        }],
                    })(
                        <Input placeholder="请输入地点名称" onChange={this.handleChange}/>
                    )}
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

export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(PlaceForm))
