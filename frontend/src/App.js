import React from 'react';
import { Routes, Route, HashRouter  } from "react-router-dom";

import Principal from './pages/Principal';

function App() {
  return (
    <HashRouter>
      <Routes>
          <Route path="/" element={<Principal/>} />
      </Routes>
    </HashRouter>  
  );
}

export default App;
