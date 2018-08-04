import React, {Component} from 'react';
import './GameSideMenu.css';
import {Segment, List, Image, Statistic, Button} from 'semantic-ui-react';
import Countdown from 'react-countdown-now';

const Completionist = () => <Statistic size='small'>
<Statistic.Label>Times up!</Statistic.Label>
</Statistic>;

export default class Timer extends Component {

  constructor(props) {
    super(props);
    this.state = {
      timeLeft: Date.now() + this.props.timer
    };
  }

  onComplete = (e) => {
    this.props.onComplete(e);
  }

  renderer = ({hours, minutes, seconds, completed}) => {
    if (completed) {
      this.onComplete();
      return <Completionist/>;
    } else {
      // Render a countdown
      return <Statistic size='small'>
        <Statistic.Value>
          {seconds}
        </Statistic.Value>
        <Statistic.Label>Seconds!</Statistic.Label>
      </Statistic>;
    }
  };

  render() {
    let date = Date.now() + 5000
    return (
      <Segment>
      <Statistic size='huge'>
        <Countdown date={this.state.timeLeft} renderer={this.renderer}>
          <Completionist />
        </Countdown>
      </Statistic>
      </Segment>
    )
  }
}
