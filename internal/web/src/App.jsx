import React, { useState, useEffect } from "react";
import "./App.css";
import Drawer from "@material-ui/core/Drawer";
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

  const [allTasks, setAllTasks] = useState([]);
  const [currentTask, setCurrentTask] = useState(null);

  useEffect(() => {
    fetch("http://localhost:8080/detail")
      .then(response => response.json())
      .then(jobs =>
        setAllTasks(
          Object.entries(jobs)
            .flatMap(([job, groups]) =>
              Object.entries(groups).flatMap(([group, tasks]) =>
                Object.entries(tasks).flatMap(([name, task]) => ({
                  ...task,
                  task: name,
                  group: group,
                  job: job,
                  vulnerable: task.trivy
                    ? task.trivy
                        .map(target => !!target.Vulnerabilities)
                        .reduce((acc, cur) => cur || acc, false)
                    : null
                }))
              )
            )
            .map((task, i) => ({ ...task, id: i }))
        )
      );
  }, []);

  return (
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
          <Sidebar
            tasks={allTasks}
            currentTask={currentTask}
            setCurrentTask={setCurrentTask}
          />
        </Drawer>
        {currentTask ? (
          <div className={classes.content}>
            <TaskBody currentTask={currentTask} />
          </div>
        ) : (
          <div className={classes.content}>
            <Typography variant="h5" color="inherit">
              Select a task
            </Typography>
          </div>
        )}
      </div>
    </ThemeProvider>
  );
}

export default App;
