/**
 * Created by Jingle on 2017/12/10.
 */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import * as _ from 'lodash'
import moment from 'moment';
import { Form, Icon, Input, Button, Select, DatePicker } from 'antd';
import { fetchData, receiveData } from '../../action';

const FormItem = Form.Item;
const Search = Input.Search;
const Option = Select.Option;

const dateFormat = 'YYYY-MM-DD';
const queryDateFormat = 'YYYYMMDD';

function hasErrors(fieldsError) {
    return Object.keys(fieldsError).some(field => fieldsError[field]);
}

class PictureSearch extends Component {

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
                let user = {}
                user.name = values.name !== undefined? values.name : ''
                user.phone = values.phone !== undefined? values.phone : ''
                const { fetchData } = this.props
                fetchData({funcName: 'searchUsers', params: user, 
                    stateName: 'usersData'})
            }
        });
    };


    getOptions = ( data=[] ) => {
        
        return data.map(item => {
            return <Option key={item.key} value={`${item.id}`}>{item.name}</Option>
        })
    }


    render() {
        const { getFieldDecorator, getFieldsError, getFieldError, isFieldTouched } = this.props.form;
        const { style, filter } = this.props
        const { townsData, countriesData, selectedDate } = this.state


        // Only show error after a field is touched.
        const fileNameError = isFieldTouched('fileName') && getFieldError('fileName');
        return (
            <Form layout="inline" style={style} onSubmit={this.handleSubmit}>
                <FormItem 
                    style={{paddingBottom: 13}}
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('name', {
                        //initialValue: townsData[0]? townsData[0].name:'',
                    })(
                        <Input
                        style={{ width: 200 }}
                        placeholder="请输入用户名"
                        >
                        </Input>
                    )}
                </FormItem>
                <FormItem
                    style={{paddingBottom: 13}}
                    validateStatus={fileNameError ? 'error' : ''}
                    help={fileNameError || ''}
                >
                    {getFieldDecorator('phone', {
                    })(
                        <Input
                        style={{ width: 200 }}
                        placeholder="请输入手机号"
                        >
                        </Input>
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

export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(PictureSearch))
