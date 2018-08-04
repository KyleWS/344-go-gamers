import React, {Component} from 'react';
import {
  Input,
  Icon,
  Grid,
  Button,
  Header,
  Loader
} from 'semantic-ui-react';
import {Divider} from 'material-ui/Divider';

export default class Welcome extends Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <Grid.Column width={10}>
        <Header as='h1'>Waiting for Other Players</Header>
        <Loader active inline='centered'/>
      </Grid.Column>
    );
  }
}