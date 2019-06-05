import React from 'react';
import styled from 'styled-components';
import ReleasesContainer from './ReleasesContainer';
import Spinner from './Spinner';

function getMostSevereAlert(events) {
  let mostSever = '';
  if (events.some(event => event.severity === 'warning')) {
    mostSever = 'warning';
  }
  if (events.some(event => event.severity === 'critical')) {
    mostSever = 'critical';
  }
  return mostSever;
}

const Header = styled.h3`
  margin: 0px;
  padding: 10px;
`;

const AppContainer = styled.div`
  background: ${props => props.backgroundColor};
  margin: 0px;  
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-bottom: 20px;
`;

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      isLoaded: false,
      events: []
    };
  }

  refetchData() {
    fetch('/api/releases')
      .then(res => res.json())
      .then(
        result => {
          let background;
          switch (getMostSevereAlert(result)) {
            case 'warning':
              background = '#ffeeba';
              break;
            case 'critical':
              background = '#E0EBF5';
              break;
            default:
              background = 'white';
          }
          this.setState({
            isLoaded: true,
            events: result,
            backgroundColor: background,
            error: null,
          });
        },
        // Note: it's important to handle errors here
        // instead of a catch() block so that we don't swallow
        // exceptions from actual bugs in components.
        error => {
          this.setState({
            isLoaded: true,
            error
          });
        }
      );
  }

  componentDidMount() {
    this.refetchData();
    this.interval = setInterval(() => this.refetchData(), 100000);
  }

  componentWillUnmount() {
    clearInterval(this.interval);
  }

  render() {
    const { error, isLoaded, events, backgroundColor } = this.state;
    if (error) {
      return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
      return <Spinner isLoading={true} />;
    } else {
      return (
        <AppContainer backgroundColor={backgroundColor}>
          <Header>Releases</Header>
          <ReleasesContainer events={events} />
        </AppContainer>
      );
    }
  }
}

export default App;
