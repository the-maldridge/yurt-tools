import React, { useState, useEffect } from "react";
import "./App.css";
import Drawer from "@material-ui/core/Drawer";
import { HashRouter as Router, Switch, Route } from "react-router-dom";
import {
  ThemeProvider,
  createMuiTheme,
  makeStyles
} from "@material-ui/core/styles";
import CssBaseline from "@material-ui/core/CssBaseline";
import green from "@material-ui/core/colors/green";
import Typography from "@material-ui/core/Typography";

import Sidebar from "./Sidebar";
import TaskBody from "./TaskBody";

const drawerWidth = 400;

const useStyles = makeStyles(theme => ({
  root: {
    display: "flex"
  },
  appBar: {
    width: `calc(100% - ${drawerWidth}px)`,
    marginLeft: drawerWidth
  },
  drawer: {
    width: drawerWidth,
    flexShrink: 0
  },
  drawerPaper: {
    width: drawerWidth
  },
  toolbar: theme.mixins.toolbar,
  content: {
    flexGrow: 1,
    width: `calc(100% - ${drawerWidth}px)`,
    padding: theme.spacing(3)
  }
}));

function App() {
  const classes = useStyles();
  const theme = React.useMemo(
    () =>
      createMuiTheme({
        palette: {
          type: true ? "dark" : "light",
          primary: green
        }
      }),
    []
  );

  const [allJobs, setAllJobs] = useState(null);

  useEffect(() => {
    fetch("/detail")
      .then(response => response.json())
      .then(jobs => setAllJobs(jobs));
  }, []);

  return (
    <Router>
      <ThemeProvider theme={theme}>
        <div className={classes.root}>
          <CssBaseline />
          <Drawer
            className={classes.drawer}
            variant="permanent"
            classes={{
              paper: classes.drawerPaper
            }}
            anchor="left"
          >
            <Sidebar jobs={allJobs} />
          </Drawer>
          <Switch>
            <Route
              path={`/task/:job/:group/:task`}
              render={({ match: { params } }) => {
                return (
                  <div className={classes.content}>
                    <TaskBody
                      currentTask={
                        allJobs
                          ? allJobs[params.job][params.group][params.task]
                          : null
                      }
                      {...params}
                    />
                  </div>
                );
              }}
            />
            <Route>
              <div className={classes.content}>
                <Typography variant="h5" color="inherit">
                  Select a task
                </Typography>
              </div>
            </Route>
          </Switch>
        </div>
      </ThemeProvider>
    </Router>
  );
}

export default App;
