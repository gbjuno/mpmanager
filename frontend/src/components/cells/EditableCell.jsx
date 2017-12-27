import React from 'react';
import { Divider, Input, Select } from 'antd';

const Option = Select.Option

const INPUT = 'input'
const SELECT = 'select'

const STRING = 'string'
const INT = 'int'

class EditableCell extends React.Component {
    state = {
        value: this.props.value,
        type: this.props.type,
        editable: this.props.editable || true,
        editType: this.props.editType || INPUT,
        valueType: this.props.valueType || STRING,
        options: this.props.options || [],
        placeholder: this.props.placeholder || '请选择对应的选项',
        style: this.props.style,
    }
    handleChange = (e) => {
        const { editType, valueType } =  this.state 
        let value = editType === SELECT ? e : e.target.value;
        switch(valueType){
            case INT:
                value = parseInt(value);
                break;
            case STRING:
                value = `${value}`
                break;
            default:
                value = `${value}`
        }
        const dataIndex = this.props.dataIndex;
        this.setState({ value });
        this.props.onChange(dataIndex, value)
    }
    handleSave = (e) => {
        //防止冒泡事件
        e.stopPropagation();
        this.setState({ editable: true });
        if (this.props.onSave) {
            this.props.onSave();
        }
    }

    handleCancel = (e) => {
        e.stopPropagation();
        this.setState({ editable: false });
        if (this.props.onCancel) {
            this.props.onCancel();
        }
    }

    edit = () => {
        this.setState({ editable: true });
    }

    getOptions = ( data=[] ) => {
        return data.map(item => {
            return <Option key={item.key} value={`${item.id}`}>{item.name}</Option>
        })
    }

    render(){
        const { value, editable, type, editType, options, placeholder, style } = this.state;

        const editUnit = ()  => {
            switch(editType){
                case 'input':
                    return (
                        <Input
                            value={value}
                            onChange={this.handleChange}
                            style={style}
                            onClick={(e) => e.stopPropagation()}
                            onPressEnter={()=>{}}
                        />
                    )
                case 'select':
                    return (
                        <Select
                            showSearch
                            defaultValue={`${value}`}
                            placeholder={placeholder}
                            onChange={this.handleChange}
                            style={{...style, width: '100%'}}
                            optionFilterProp="children"
                            filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                        >
                            {this.getOptions(options)}
                        </Select>
                    )
                default:
                    return <p>Unkown Edit Type</p>
            }
        }

        return (
            <div className="editable-cell">
                {
                editable ?
                    type !== "opt"?
                        <div className="editable-cell-input-wrapper">
                        {editUnit()}
                        </div>
                        :
                        <span>
                            <a className="opt-confirm" onClick={this.handleSave}>保存</a>
                            <Divider type="vertical" />
                            <a className="opt-cancel"  onClick={this.handleCancel}>取消</a>
                        </span>
                    :
                    <a target="_blank">{value}</a>
                }
            </div>
        )
    }
}

export default EditableCell
