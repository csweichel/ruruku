import * as React from 'react';
import * as ReactDOM from 'react-dom';
import App from './App';
import { MiniEventEmitter } from './types/mini-event-emitter';
import './index.css';
import registerServiceWorker from './registerServiceWorker';
import * as showdown from 'showdown';

const reloadRequest = new MiniEventEmitter<boolean>();

showdown.setOption('simplifiedAutoLink', true);

ReactDOM.render(
  <App reloadRequest={reloadRequest} />,
  document.getElementById('root') as HTMLElement
);
registerServiceWorker(() => reloadRequest.publish(true));
