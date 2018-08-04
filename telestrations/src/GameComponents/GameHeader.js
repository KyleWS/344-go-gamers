import React, {Component} from 'react';
import {
    Grid,
    Container,
    Header,
    Icon,
    Segment
} from 'semantic-ui-react';

export default class GameHeader extends Component {
    render() {
        return (
            <Segment basic>
            <span>
            
            <Header as='h2' icon textAlign='center'>
                <span>
                <Icon name='hourglass end' circular />
                    <Header.Content>
                        
                        Round 3
                    </Header.Content>
            </span>
                </Header>
            </span>
            </Segment>
        );
    }
}