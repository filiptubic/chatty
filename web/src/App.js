import React from "react";
import { Box } from '@mui/material';
import './App.css';
import Chat from './components/Chat';
import Header from './components/Header';
import { Sidebar, DrawerHeader } from './components/Sidebar'
import {
  BrowserRouter,
  Routes,
  Route,
} from "react-router-dom";

function App() {

  const [open, setOpen] = React.useState(false);

  const handleDrawerClose = () => {
    setOpen(false);
  };

  return (
    <div className="App">
      <DrawerHeader />
      <BrowserRouter>
        <Box sx={{ display: 'flex' }}>
          <Header
            open={open}
            setOpen={setOpen}
            handleDrawerClose={handleDrawerClose}
            sx={{ flexGrow: 1 }}
          />
          <Sidebar
            sx={{ flexGrow: 1 }}
            open={open}
            setOpen={setOpen}
            handleDrawerClose={handleDrawerClose}
          />
          <Box sx={{ flexGrow: 1 }}>
            <Routes>
              <Route path="/:chatId" index element={ <Chat /> }/>
              <Route path="/" element={<div></div>} />
            </Routes>
          </Box>
        </Box>
      </BrowserRouter>
    </div>
  );
}

export default App;
