import React, {Component} from 'react';
import './GameSideMenu.css';
import {Segment, List, Image, Statistic, Button} from 'semantic-ui-react';

export default class SideMenu extends Component {
  constructor(props) {
    super(props);
  }

  buildPlayers() {
    let players = this.props.players;
    const playerItems = [];
    if (players.length !== 0) {
      (players).forEach(player => {
        playerItems.push(
          <List.Item>
          <Image avatar src={player.photoURL}/>
            <List.Content className="list-content">
              <List.Header >{player.firstName + " " + player.lastName}</List.Header>
              <div className="description">{player.userName}</div>
            </List.Content>
          </List.Item>
        )
      });
    } else {
      playerItems.push(<h4 key={"waiting"}>Waiting for players...</h4>)
    }
    return playerItems;
  }

  render() {
    let playerItems = this.buildPlayers();
    return (
      <div>
        <Segment>
          <Statistic size='huge'>
            <Statistic.Label>Views</Statistic.Label>
          </Statistic>
          <h2>Other Players:</h2>
          <List relaxed size='large'>
            {playerItems}
          </List>
        </Segment>
      </div>
    );
  }
}