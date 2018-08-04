import React, {Component} from 'react';

import './App.css';

import {Grid} from 'semantic-ui-react';
import axios from 'axios';
import {CONST} from './Constants/Constants';
import SockJS from 'sockjs-client'
import LoginScreen from './loginScreen';
import {
  GameHeader,
  Menu,
  Canvas,
  TextInput,
  SideMenu,
  StartButton,
  Timer,
  Results,
  Welcome
} from './GameComponents/GameComponents.js'
import Countdown from 'react-countdown-now';
/*
Module:superagent
superagent is used to handle post/get requests to server
*/

export default class Game extends Component {

  constructor(props) {
    super(props);
    this.authHeader = window.localStorage.auth;    
    this.sock = new WebSocket('wss://api.telestrations.alexsirr.me/v1/game/ws?auth=' + this.authHeader);
    this.img64 = "";
    this.state = {
      base: null,      
      timer: null,
      display: [< Welcome />],
      toDraw: '',
      toDesc: '',
      currGuess: '',
      currSketch: '{"objects":""}',
      results: [],
      seconds: 0,
      actions: this.sock,
      game: {
        user: window.localStorage.user,
        players: []
      },
      phase:'start',
    }
  }
  componentDidMount = () =>{
    this.saveImg = this.saveImg.bind(this);
    this.sock.onopen = () => {
      this.sock.send("first send");

    }
    this.sock.onmessage = e => {
      //PARSE JSON
      var res = JSON.parse(e.data);
      //SWITCH ON ACTIION
      switch (res.action) {
        case 'game-start':
          this.setState({
            game: {
              players: res.players
            }
          });
          break;
        case 'first-phase':
          this.setState({
            phase: 'desc',
            display: [< TextInput intro = {
                "Choose something dope as shit, boring examples picked up by our super gangsta AI" +
                  " will be kicked"
              }
              newMsg = {
                this.onTextInput
              } />],
            seconds: res["round-duration"]
          });
          // Calls buildTimer method that will build and start the timer object
          this.buildTimer();
          break;
        case 'drawing-phase':
          this.buildTimer();
          this.setState({toDraw: res.data,
          phase: "draw",seconds: res["round-duration"]});
          this.setState({
            toDraw: res.data,
            phase: "draw",
            display: [< Canvas toDraw = {
                this.state.toDraw
              }
              saveImg = {
                this.saveImg 
              } auth={this.authHeader} save={this.saveImg}/>]
          });
          this.buildTimer();        
          break;
        case 'description-phase':
      this.save64(res.data).then(response => {
        this.setState({base: JSON.parse(JSON.stringify(response.data.data))});   
        this.setState({
          phase: 'desc',
          toDesc: res.data,
          display: [< TextInput intro = {
              "Describe the deeper meaning behind this Bob Ross masterpiece"
            }
            toDesc = {
              this.state.toDesc
            }
            newMsg = {
              this.onTextInput
            }
            base = {this.state.base}
             />],
          seconds: res["round-duration"]
        });
        this.buildTimer();     
      });
          break;
        case 'game-results':
          this.setState({
            results: res.results,
            display: [<Welcome />],
            seconds: res["round-duration"]
          });
          break;
        case 'intermission':
          this.setState({
            display: [< Welcome />],
            seconds: res["duration"]
          });
          break;
        default:
          return;
      }
    }
  }
  onTextInput = (g) => {
    this.setState({currGuess: g});

  }
  saveImg = (img) => {
    this.setState({currSketch: img});
  }
  //utilizes return axios for get and response handling
  save64 = (res) =>{
    return axios({
      method: 'GET',
      url: CONST.API_URL + "/static/" + res,
      responseType: 'json',
        headers: {
          'Authorization': this.authHeader
        }
      });
  }
  sendDescription = () => {
    var descJSON = {
      type: "description",
      data: this.state.currGuess
    }
    descJSON = JSON.stringify(descJSON);
    axios({
      method: 'POST',
      url: CONST.API_URL + "/v1/game/submit",
      data: JSON.parse(descJSON),
        headers: {
          'Authorization': this.authHeader
        }
      })
      .then(function (response) {
        if (response.status == 200) {
        } else {
          console.log('did not send');
        }
      })
      .catch(function (error) {
        console.log(error);
      });
      this.setState({currGuess: ''});
  }
  buildTimer = () => {
    this.setState({timer: null});
    this.setState({timer: <Timer timer={this.state.seconds * 1000} onComplete={this.handleTimerComplete} phase={this.state.phase}/>})
  }
 
  handleTimerComplete = () => {
    if(this.state.phase == 'desc'){
    this.sendDescription();    
      }
      if (this.state.phase =='draw') {
              var imgJSON = {
                type: "drawing",
                data: this.state.currSketch
              }
              imgJSON = JSON.stringify(imgJSON);
              axios({
                method: 'POST',
                url: CONST.API_URL + "/v1/game/submit",
                data: JSON.parse(imgJSON),
                  headers: {
                    'Authorization': this.authHeader
                  }
                })
                .then(function (response) {
                  if (response.status == 200) {
                    console.log("sent");
                  } else {
                    console.log('did not send');
                  }
                })
                .catch(function (error) {
                  console.log(error);
                });
                this.setState({currSketch: ''});  
      }
  }
  render() {
    var props = this.state.game;
    var gameStyle = {
      padding: '1.5em'
    }
    return (
      <div >
        <Menu {...this.props} user={this.state.game.user}/>
        <Grid className="game" style={gameStyle}>
          <Grid.Column width={3}>
              {this.state.timer}
            <SideMenu {...props}/>
          </Grid.Column>
          {this.state.display}
        </Grid>
      </div>
    );
  }
}

// WEBPACK FOOTER // src/Game.js