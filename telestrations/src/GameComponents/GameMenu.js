import React, {Component} from 'react'
import {
  Menu,
  Button,
  Modal,
  Image,
  Header,
  Icon
} from 'semantic-ui-react'
import LoginScreen from '../loginScreen'

export default class GameMenu extends Component {
  constructor(props) {
    super(props);
    this.state = {
      activeItem: 'home',
      user: JSON.parse(this.props.user)
    }
  }

  handleItemClick = (e, {name}) => this.setState({activeItem: name});

  handleLogout = () => {
    var self = this;
    var loginPage = [];
    loginPage.push(<LoginScreen appContext={this.props.appContext}/>);
    console.log(self.props);
    this.props.appContext.setState({loginPage: loginPage, gameScreen: []})
  }

  componentWillMount() {

  }

  render() {
    const {activeItem} = this.state
    return (
      <Menu size='large'>
        <Menu.Item
          name='Home'
          active={activeItem === 'Home'}
          onClick={this.handleItemClick}/>
        <Modal
          trigger={< Menu.Item name = 'Profile' active = {
          activeItem === 'Profile'
        }
        onClick = {this.handleItemClick} />}>
          <Modal.Content>
            <Header as='h2' icon textAlign='center'>
              <Icon name='users' circular/>
              <Header.Content>
                {this.state.user.firstName + " " + this.state.user.lastName}
              </Header.Content>
            </Header>
            {"Email: " + this.state.user.email}
            <br/>
            {"UserName: " + this.state.user.userName}
          </Modal.Content>
        </Modal>
        <Menu.Menu style={{ marginLeft:"23%", paddingTop: "10px", fontVariant: "small-caps" }}>
        <h3>
         Digital Telestrations
         <Icon name="pencil" style={{marginLeft: "10px" }} color="blue"/>
        </h3>
        </Menu.Menu>


        <Menu.Menu position='right'>
          <Button color='red' onClick={this.handleLogout}>
            Logout
          </Button>
        </Menu.Menu>
      </Menu>
    )
  }
}