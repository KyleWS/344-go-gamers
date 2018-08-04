import React, { Component } from 'react'
import { Grid, Icon, Segment, Dropdown, Header, Button, Divider } from 'semantic-ui-react'
import GameHeader from './GameHeader';
import {CONST} from '../Constants/Constants';
import axios from 'axios';
import { SketchField, Tools } from 'react-sketch';

import {
    Card,
    CardHeader,
    CardText,
    GridList,
    GridTile,
    IconButton,
    MenuItem,
    Slider,
    Toggle,
} from 'material-ui';
import { CompactPicker } from 'react-color';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import RemoveIcon from 'material-ui/svg-icons/content/clear';
import injectTapEventPlugin from 'react-tap-event-plugin';
injectTapEventPlugin();

const styles = {
    root: {
        padding: '3px',
        display: 'flex',
        flexWrap: 'wrap',
        margin: '10px 10px 5px 10px',
        justifyContent: 'space-around'
    },
    gridList: {
        width: '100%',
        overflowY: 'auto',
        marginBottom: '0px'
    },
    gridTile: {
        backgroundColor: '#fcfcfc'
    },
    separator: {
        height: '14px',
        backgroundColor: 'white'
    },
    iconButton: {
        fill: 'white',
        width: '42px',
        height: '42px'
    },
    dropArea: {
        width: '100%',
        height: '64px',
        border: '2px dashed rgb(102, 102, 102)',
        borderStyle: 'dashed',
        borderRadius: '5px',
        textAlign: 'center',
        paddingTop: '20px'
    },
    activeStyle: {
        borderStyle: 'solid',
        backgroundColor: '#eee'
    },
    rejectStyle: {
        borderStyle: 'solid',
        backgroundColor: '#ffdddd'
    }
};


/**
* Helper function to manually fire an event
*
* @param el the element
* @param etype the event type
*/
function eventFire(el, etype) {
    if (el.fireEvent) {
        el.fireEvent('on' + etype);
    } else {
        var evObj = document.createEvent('Events');
        evObj.initEvent(etype, true, false);
        el.dispatchEvent(evObj);
    }
}

export default class Canvas extends Component {
    constructor(props) {
        super(props);
        this.state = {
            lineColor: 'black',
            lineWidth: 10,
            fillColor: '#68CCCA',
            backgroundColor: 'transparent',
            shadowWidth: 0,
            shadowOffset: 0,
            tool: Tools.Pencil,
            fillWithColor: false,
            fillWithBackgroundColor: false,
            drawings: [],
            // canUndo: false,
            // canRedo: false,
            controlledSize: false,
            sketchWidth: 600,
            sketchHeight: 600,
            stretched: true,
            stretchedX: false,
            stretchedY: false,
            originX: 'left',
            originY: 'top',
            desc: this.props.toDraw
        };
    }

    _selectTool = (event, data) => {
        this.setState({
            tool: data.value
        });
    };
    _save = () => {
        if(this._sketch != null){
       this.props.save(this._sketch.toDataURL());
        }
    };
    _download = () => {
        /*eslint-disable no-console*/

        console.save(this._sketch.toDataURL(), 'toDataURL.txt');
        console.save(JSON.stringify(this._sketch.toJSON()), 'toDataJSON.txt');

        /*eslint-enable no-console*/

        let { imgDown } = this.refs;
        let event = new Event('click', {});

        imgDown.href = this._sketch.toDataURL();
        imgDown.download = 'toPNG.png';
        imgDown.dispatchEvent(event);
    };
    _renderTile = (drawing, index) => {
        return (
            <GridTile
                key={index}
                title='Canvas Image'
                actionPosition="left"
                titlePosition="top"
                titleBackground="linear-gradient(to bottom, rgba(0,0,0,0.7) 0%,rgba(0,0,0,0.3) 70%,rgba(0,0,0,0) 100%)"
                cols={1} rows={1} style={styles.gridTile}
                actionIcon={<IconButton onTouchTap={(c) => this._removeMe(index)}><RemoveIcon
                    color="white" /></IconButton>}>
                <img src={drawing} />
            </GridTile>
        );
    };
    _removeMe = (index) => {
        let drawings = this.state.drawings;
        drawings.splice(index, 1);
        this.setState({ drawings: drawings });
    };
    // _undo = () => {
    //     this._sketch.undo();
    //     this.setState({
    //         canUndo: this._sketch.canUndo(),
    //         canRedo: this._sketch.canRedo()
    //     })
    // };
    // _redo = () => {
    //     this._sketch.redo();
    //     this.setState({
    //         canUndo: this._sketch.canUndo(),
    //         canRedo: this._sketch.canRedo()
    //     })
    // };
    _clear = () => {
        this._sketch.clear();
        this._sketch.setBackgroundFromDataUrl('');
        this.setState({
            controlledValue: null,
            backgroundColor: 'transparent',
            fillWithBackgroundColor: false,
            // canUndo: this._sketch.canUndo(),
            // canRedo: this._sketch.canRedo()
        })
    };
    _onSketchChange = () => {
        this._save();
        // let prev = this.state.canUndo;
        // let now = this._sketch.canUndo();
        // if (prev !== now) {
        //     this.setState({ canUndo: now });
        // }
    };
    _onBackgroundImageDrop = (accepted/*, rejected*/) => {
        if (accepted && accepted.length > 0) {
            let sketch = this._sketch;
            let reader = new FileReader();
            let { stretched, stretchedX, stretchedY, originX, originY } = this.state;
            reader.addEventListener('load', () => sketch.setBackgroundFromDataUrl(reader.result, {
                stretched: stretched,
                stretchedX: stretchedX,
                stretchedY: stretchedY,
                originX: originX,
                originY: originY
            }), false);
            reader.readAsDataURL(accepted[0]);
        }
    };

    _selectColor = () => { }
    componentDidMount = () => {

        /*eslint-disable no-console*/

        (function (console) {
            console.save = function (data, filename) {
                if (!data) {
                    console.error('Console.save: No data');
                    return;
                }
                if (!filename) filename = 'console.json';
                if (typeof data === 'object') {
                    data = JSON.stringify(data, undefined, 4)
                }
                var blob = new Blob([data], { type: 'text/json' }),
                    e = document.createEvent('MouseEvents'),
                    a = document.createElement('a');
                a.download = filename;
                a.href = window.URL.createObjectURL(blob);
                a.dataset.downloadurl = ['text/json', a.download, a.href].join(':');
                e.initMouseEvent('click', true, false, window, 0, 0, 0, 0, 0, false, false, false, false, 0, null);
                a.dispatchEvent(e)
            }
        })(console);
    };
    render() {
        let colorOptions = [
            <CompactPicker
                id='lineColor' color={this.state.lineColor}
                onChange={(color) => this.setState({ lineColor: color.hex })} />

        ];

        let fillColorOptions = [
            <Toggle label="Fill"
                defaultToggled={this.state.fillWithColor}
                onToggle={(e) => this.setState({ fillWithColor: !this.state.fillWithColor })} />,
            <CompactPicker
                color={this.state.fillColor}
                onChange={(color) => this.setState({ fillColor: color.hex })} />
        ];
        let toolOptions = [
            {
                text: 'Select',
                value: Tools.Select,
                image: <Icon name='mouse pointer' />
            },
            {
                text: 'Pencil',
                value: Tools.Pencil,
                image: <Icon name='pencil' />
            },
            {
                text: 'Line',
                value: Tools.Line,
                image: <Icon name='minus' />
            },
            {
                text: 'Rectangle',
                value: Tools.Rectangle,
                image: <Icon name='square' />
            },
            {
                text: 'Circle',
                value: Tools.Circle,
                image: <Icon name='circle' />
            }
        ];

        let { controlledValue } = this.state;
        return (
            <MuiThemeProvider muiTheme={getMuiTheme()}>
                <Grid.Column width={10}>
                    <Segment>
                        <SketchField
                            name='sketch'
                            className='canvas-area'
                            ref={(c) => this._sketch = c}
                            lineColor={this.state.lineColor}
                            lineWidth={this.state.lineWidth}
                            fillColor={this.state.fillWithColor ? this.state.fillColor : 'transparent'}
                            backgroundColor={this.state.fillWithBackgroundColor ? this.state.backgroundColor : 'transparent'}
                            width={this.state.controlledSize ? this.state.sketchWidth : null}
                            height={this.state.controlledSize ? this.state.sketchHeight : null}
                            // defaultValue={dataJson}
                            value={controlledValue}
                            forceValue={true}
                            onChange={this._onSketchChange}
                            tool={this.state.tool} />
                    </Segment>
                </Grid.Column>
                <Grid.Column width={3}>
                    <Segment>
                        <Button primary name="submit" onClick={this._download}>{this.props.toDraw}</Button>
                    </Segment>

                    <Segment>
                        {/* <Icon link size="big"
                            name="undo"
                            onClick={this._undo}
                            iconStyle={styles.iconButton}
                            disabled={!this.state.canUndo}>
                        </Icon>
                        <Icon link size="big"
                            name="undo"
                            flipped="horizontally"
                            onClick={this._redo}
                            iconStyle={styles.iconButton}
                            disabled={!this.state.canRedo}>
                        </Icon> */}
                        <a ref='imgDown' />
                        <Divider hidden />

                        <div className='row'>

                            <Header as='h5'>  Drawing Tool </Header>
                            <Dropdown
                                placeholder='Select Tool' fluid selection
                                options={toolOptions}
                                onChange={this._selectTool} />
                            <Header as='h5'>Line Color</Header>
                            <Dropdown
                                placeholder='Line Color' fluid selection
                                text={this.state.lineColor}
                                options={colorOptions}
                                onChange={(color) => this.setState({ lineColor: color.hex })} />
                            <Header as='h5'>Fill Color</Header>
                            <Dropdown
                                placeholder='Fill Color' fluid selection
                                options={fillColorOptions}
                                color={this.state.fillColor}
                                onChange={(color) => this.setState({ fillColor: color.hex })} />
                            <Header as='h5'>  Brush Width </Header>
                            <Slider ref='slider' step={0.1}
                                defaultValue={this.state.lineWidth / 100}
                                onChange={(e, v) => this.setState({ lineWidth: v * 100 })} />

                            <Card>
                                <CardHeader title='Colors' actAsExpander={true} showExpandableButton={true} />
                                <CardText expandable={true}>
                                    <label htmlFor='lineColor'>Line</label>
                                    <CompactPicker
                                        id='lineColor' color={this.state.lineColor}
                                        onChange={(color) => this.setState({ lineColor: color.hex })} />
                                    {/* <Toggle label="Fill"
                                      defaultToggled={this.state.fillWithColor}
                                      onToggle={(e) => this.setState({fillWithColor: !this.state.fillWithColor})}/>
                              <CompactPicker
                                  color={this.state.fillColor}
                                  onChange={(color) => this.setState({fillColor: color.hex})}/> */}
                                </CardText>
                            </Card>
                        </div>

                        <div className='row'>
                            <div className="col-xs-12 col-sm-12 col-md-12 col-lg-12">
                                <div className="box" style={styles.root}>
                                    <GridList
                                        cols={5}
                                        cellHeight={200}
                                        padding={1} style={styles.gridList}>
                                        {this.state.drawings.map(this._renderTile)}
                                    </GridList>
                                </div>
                            </div>
                        </div>
                    </Segment>
                </Grid.Column>
            </MuiThemeProvider>

        )
    }
}
