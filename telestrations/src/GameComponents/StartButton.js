import React, { Component } from 'react';
import { Segment, Button } from 'semantic-ui-react';

export default class SideMenu extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
        <Segment>
          {/* Calls parent function in Game.js to start the game */}
          <Button color='green' name="start" onClick={this.props.myClick}>Start Game!</Button>
        </Segment>
    );
  }
}