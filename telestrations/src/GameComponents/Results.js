import React, {Component} from 'react';
import {
  Input,
  Icon,
  Grid,
  Button,
  List,
  Image
} from 'semantic-ui-react';
import axios from 'axios';
import {CONST} from '../Constants/Constants';

export default class Results extends Component {
  constructor(props) {
    super(props);
        this.results = {"action":"game-results","results":[[{"owner":{"id":"5a273fb4b2b95b0001fafd37","email":"calvin.korver@gmail.com","userName":"cjkorver","firstName":"Calvin","lastName":"Korver","photoURL":"https://www.gravatar.com/avatar/d1a1a4f8247e7c8cec0b753686d6b4a1"},"imageHash":"","description":"round2"},{"owner":{"id":"5a2984e923fd480001d4466e","email":"bob@gmail.com","userName":"bob","firstName":"bob","lastName":"sparks","photoURL":"https://www.gravatar.com/avatar/0ca52abaa9fce14c44a351fccd1b9fc5"},"imageHash":"665038abce22e8f5b5f590ef1bf8e626d5dbd8a9ef86d62f2018e62bf077ca64","description":""},{"owner":{"id":"5a273fb4b2b95b0001fafd37","email":"calvin.korver@gmail.com","userName":"cjkorver","firstName":"Calvin","lastName":"Korver","photoURL":"https://www.gravatar.com/avatar/d1a1a4f8247e7c8cec0b753686d6b4a1"},"imageHash":"","description":"rond3"}],[{"owner":{"id":"5a2984e923fd480001d4466e","email":"bob@gmail.com","userName":"bob","firstName":"bob","lastName":"sparks","photoURL":"https://www.gravatar.com/avatar/0ca52abaa9fce14c44a351fccd1b9fc5"},"imageHash":"","description":"round2"},{"owner":{"id":"5a273fb4b2b95b0001fafd37","email":"calvin.korver@gmail.com","userName":"cjkorver","firstName":"Calvin","lastName":"Korver","photoURL":"https://www.gravatar.com/avatar/d1a1a4f8247e7c8cec0b753686d6b4a1"},"imageHash":"01681664055cd83b966b29cfd6bb9f43f163d6357cc0b86e54e600a165299e2a","description":""},{"owner":{"id":"5a2984e923fd480001d4466e","email":"bob@gmail.com","userName":"bob","firstName":"bob","lastName":"sparks","photoURL":"https://www.gravatar.com/avatar/0ca52abaa9fce14c44a351fccd1b9fc5"},"imageHash":"","description":"fuck"}]]}
    }


  buildResults = (results) => {

    let resultsFormatted = [];
    results.results.forEach((roundArray, index) => {
      console.log("Round #" + index);
      roundArray.forEach(playerElement => {
        if (playerElement.imageHash.length !== 0) {
          let image = this.getImage(playerElement.imageHash);
          console.log(image);
        }
      });
    });
  }

  getImage = (hashedURL) => {
    return axios({
      method: 'GET',
      url: CONST.API_URL + "/static?" + hashedURL,
      responseType: 'json',
      data: hashedURL,
        headers: {
          'Authorization': window.localStorage.auth
        }
      })
  }

  //   let playerRes = [];
  //   if (results !== undefined) {
  //     playerRes = results.map((result) => <List.Item>
  //       {/* need to make this string back into image */}
  //       <Image src={result.ImageHash}/>
  //       <div class="description">{result.Description}</div>
  //       <List.Content className="list-content">
  //         <List.Header >{result.userStruct.firstName + " " + result.userStruct.lastName}</List.Header>
  //         <div class="description">{result.userStruct.userName}</div>
  //       </List.Content>
  //     </List.Item>);
  //   }
  //   return playerRes;
  // }

  render() {

    let playerRes = this.buildResults(this.results);

    return (
      <Grid.Column width={10}>
        <List relaxed size='large'>
          {playerRes}
        </List>
      </Grid.Column>
    );
  }
}