import React, {Component} from 'react';
import {Input, Icon, Grid, Button, Header, Image} from 'semantic-ui-react';
import {Divider} from 'material-ui/Divider';

export default class TextInput extends Component {
  constructor(props) {
    console.log(props);
    super(props);
    this.state = {
      guess: "",
      desc: this.props.toDesc,
    }
  }

  handleMessage (e) { 
    this.setState({guess: e.target.value});
    this.props.newMsg(e.target.value);
  };
  render() {
    return (
        <Grid.Column width={10}>
        <Header>{this.props.intro}</Header>
        <Image src={this.props.base} />
            <Input width={6} type = "text" value={this.state.guess} onChange={this.handleMessage.bind(this)} placeholder='Make a guess!'/>
        </Grid.Column>
    );
  }
}