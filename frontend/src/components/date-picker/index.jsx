import React from 'react'
import moment from 'moment';
import { DatePicker } from 'antd';

const modeValue={
    date: 0,
    month: 1,
    year: 2,
    decade: 3
};

export default class MyDatePicker extends React.Component {
    static defaultProps={
        topMode: "month",
        format: "YYYY-MM-DD"
    };
    constructor(props) { 
        super(props);        
        this.state={
            value: this.props.value || this.props.defaultValue,
            mode: this.props.topMode,
            preMode: this.props.topMode
        };
        this.isOnChange = false;
    }
    componentWillReceiveProps(nextProps, nextContext){
        if(this.props.topMode != nextProps.topMode){
            this.setState({
                mode: nextProps.topMode
            });
        }
    }
    /**
     * 
     * @param {*} value 
     * @param {*} mode 
     */
    onPanelChange(value, mode){
        // console.log(`onPanelChange date:${value} mode:${mode}`);
        mode = mode || "month";
        let open = true;
        if(modeValue[this.state.mode] > modeValue[mode] && modeValue[this.props.topMode] > modeValue[mode]) {
            //向下
            open = false;
            mode = this.props.topMode;
        }
        this.setState({
            value, open, mode,
            preMode: this.state.mode
        });

        if(this.props.onChange){
            this.props.onChange(value)
        }
    }

    onChange(value, dateStr){
        // console.log(`onChange date:${value} dateStr:${dateStr}`);
        this.isOnChange = true;
        this.setState({
            open: false,
            value
        });
        if(this.props.onChange){
            this.props.onChange(value)
        }
    }
    
    render() {
        return (
            <DatePicker 
                allowClear
                value={this.state.value} 
                mode={this.state.mode}
                open={this.state.open}
                format={this.props.format}
                onFocus={()=>!this.isOnChange&&(this.isOnChange=!this.isOnChange,this.setState({open:true}))}
                onChange={this.onChange.bind(this)}
                onPanelChange={this.onPanelChange.bind(this)}
                onOpenChange={(open)=>this.setState({open})}
            />
        );
    }

}