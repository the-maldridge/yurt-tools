import React, { useState } from "react";
import List from "@material-ui/core/List";
import Divider from "@material-ui/core/Divider";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import ListItemText from "@material-ui/core/ListItemText";
import TextField from "@material-ui/core/TextField";
import green from "@material-ui/core/colors/green";
import red from "@material-ui/core/colors/red";
import LockIcon from "@material-ui/icons/Lock";
import UpdateIcon from "@material-ui/icons/Update";

export default function Sidebar({ tasks, currentTask, setCurrentTask }) {
  const [onlySec, setOnlySec] = useState(false);
  const [onlyUpdate, setOnlyUpdate] = useState(false);
  const [searchText, setSearchText] = useState("");

  return (
    <List component="nav">
      <ListItem button>
        <ListItemIcon>
          <LockIcon
            style={onlySec ? { color: red[200] } : {}}
            onClick={() => setOnlySec(!onlySec)}
          />
        </ListItemIcon>
        <ListItemIcon>
          <UpdateIcon
            style={onlyUpdate ? { color: red[200] } : {}}
            onClick={() => setOnlyUpdate(!onlyUpdate)}
          />
        </ListItemIcon>
        <ListItemText>
          <TextField
            variant="outlined"
            size="small"
            label="Search"
            onChange={event => setSearchText(event.target.value)}
          />
        </ListItemText>
      </ListItem>
      <Divider />
      {tasks
        .filter(
          task =>
            task.task.includes(searchText) ||
            task.group.includes(searchText) ||
            task.job.includes(searchText)
        )
        .filter(task => !onlySec || task.vulnerable)
        .filter(
          task => !onlyUpdate || (task.versions && !task.versions.UpToDate)
        )
        .map(task => (
          <ListItem
            button
            selected={currentTask && currentTask.id === task.id}
            key={task.id}
            onClick={() => {
              setCurrentTask(task);
            }}
          >
            <ListItemIcon>
              {task.vulnerable != null ? (
                task.vulnerable ? (
                  <LockIcon style={{ color: red[200] }} />
                ) : (
                  <LockIcon style={{ color: green[200] }} />
                )
              ) : (
                <LockIcon />
              )}
            </ListItemIcon>
            <ListItemIcon>
              {task.versions ? (
                task.versions.UpToDate ? (
                  <UpdateIcon style={{ color: green[200] }} />
                ) : (
                  <UpdateIcon style={{ color: red[200] }} />
                )
              ) : (
                <UpdateIcon />
              )}
            </ListItemIcon>
            <ListItemText
              primary={`${task.task}`}
              secondary={`${task.job} » ${task.group}`}
            />
          </ListItem>
        ))}
    </List>
  );
}
