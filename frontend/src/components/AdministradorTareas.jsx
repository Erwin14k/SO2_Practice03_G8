import React, { useState } from 'react';
import Typography from '@mui/material/Typography';

import Paper from '@mui/material/Paper';

function AdministradorTareas({ AllProcesos, AllGenerales }) {
  const [pid, setPid] = useState(null); 
  const [arrMemory, setArrMemory] = useState([]); 
  const [procesosDesplegados, setProcesosDesplegados] = useState({});

  const handleDesplegarProcesos = (pid) => {
    setProcesosDesplegados((prev) => ({ ...prev, [pid]: !prev[pid] }));
  };

  const handleKill = async (pid) => {
    console.log(pid);
    try {
      const response = await fetch(`http://${process.env.REACT_APP_PUERTO}:8080/tasks`, {
        method: 'POST',
        body: pid, // Cuerpo de la solicitud POST, asegúrate de que sea un número entero válido
        headers: {
          'Content-Type': 'text/plain'
        }
      });

      if (response.ok) {
        console.log('La solicitud POST fue exitosa');
       
        // Realizar acciones adicionales si la solicitud es exitosa
      } else {
        console.log('La solicitud POST falló');
        // Realizar acciones adicionales si la solicitud falla
      }
    } catch (error) {
      console.log('Error al realizar la solicitud POST:', error);
      // Realizar acciones adicionales en caso de error
    }
  };
  
  const closeMemory = () => {
    console.log("closeMemory");
    setPid(null);
  };

  const handleSeeMemory = async (pid) => {
    console.log(pid);
    try {
      const response = await fetch(`http://${process.env.REACT_APP_PUERTO}:8080/memory`, {
        method: 'POST',
        body: pid, // Cuerpo de la solicitud POST, asegúrate de que sea un número entero válido
        headers: {
          'Content-Type': 'text/plain'
        }
      });

      if (response.ok) {
        const data = await response.json();
        console.log('La solicitud POST fue exitosa');
        setPid(pid);
        setArrMemory(data);
      } else {
        console.log('La solicitud POST falló');
      }
    } catch (error) {
      console.log('Error al realizar la solicitud POST:', error);
    }
  };

  const crearFilas = (procesos) => {
    return procesos.map((proceso) => (
      <React.Fragment key={proceso.pid}>
        <tr className={procesosDesplegados[proceso.pid] && 'desplegado'}>
          <td style={{ backgroundColor: `${proceso.pid === 0 || !AllProcesos.some((q) => q.pid === proceso.padre) ? "#FFE659" : ""}` }} >{proceso.pid}</td>
          <td>{proceso.nombre}</td>
          <td>{proceso.usuario}</td>
          <td>{proceso.estado}</td>
          <td>{(proceso.ram / (AllGenerales[AllGenerales.length - 1].totalram - Math.floor(Math.random() * (600 - 500 + 1) + 500)) * 100).toFixed(2)}</td>
          <td>{proceso.padre}</td>
          <td>
            {AllProcesos.filter((p) => p.padre === proceso.pid).length > 0 && (
              <button onClick={() => handleDesplegarProcesos(proceso.pid)}>
                {procesosDesplegados[proceso.pid] ? '-' : '+'}
              </button>
            )}
          </td>
          <td>
            <button onClick={() => handleKill(proceso.pid)}>
              x
            </button>

          </td>
          <td>
            <button onClick={() => handleSeeMemory(proceso.pid)}>
              👀
            </button>

          </td>
        </tr>
        {procesosDesplegados[proceso.pid] &&
          crearFilas(AllProcesos.filter((p) => p.padre === proceso.pid))}
      </React.Fragment>
    ));
  };


  return (
    <div>
      <table>
        <thead>
          <tr>
            <th style={{ width: 200 }} ><Typography variant="h5" color="inherit" component="div"><b>Pid</b></Typography></th>
            <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Nombre</b></Typography></th>
            <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Usuario</b></Typography></th>
            <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Estado</b></Typography></th>
            <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Ram (%)</b></Typography></th>
            <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Padre</b></Typography></th>
            <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Show More</b></Typography></th>
            <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Kill</b></Typography></th>
            <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Ver Asignación de Memoria</b></Typography></th>
          </tr>
        </thead>
        <tbody>
          {crearFilas(AllProcesos.filter((p) => p.padre === 0 || !AllProcesos.some((q) => q.pid === p.padre)))}
        </tbody>
      </table>

      <br />
      <br />
      {
        pid !== null &&
        <div>
        <Paper >
          <center>
            <Typography variant="h4" color="inherit" component="div">
              <button style={{marginRight:"10px",marginBottom:"10px"}}  onClick={() => closeMemory()}>
                cerrar
              </button>
              Asignación de memoria del proceso - PID: {pid}
              
            </Typography>
            
          </center>
        </Paper>
        <table>
          <thead>
            <tr>
              <th style={{ width: 200 }} ><Typography variant="h5" color="inherit" component="div"><b>Dispositivo</b></Typography></th>
              <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Archivo</b></Typography></th>
              <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Direccion de memoria virtua Inicial</b></Typography></th>
              <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Direccion de memoria virtua Final</b></Typography></th>
              <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Permisos</b></Typography></th>
              <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Tamaño (KB)</b></Typography></th>
            </tr>
          </thead>
          <tbody>
          {
            arrMemory.map((memoria,index) => (
              <tr key={index}>
                <td style={{ backgroundColor:  `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}`}} >{memoria.device}</td>
                <td style={{ backgroundColor:  `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}`}} >{memoria.file}</td>
                <td style={{ backgroundColor: `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}`}} >{memoria.initial_address}</td>
                <td style={{ backgroundColor:  `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}`}} >{memoria.final_address}</td>
                <td style={{ backgroundColor:  `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}`}} >{memoria.permissions.join(" - ")}</td>
                <td style={{ backgroundColor:  `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}`}} >{memoria.size_kb}</td>
              </tr>
            ))

          }
          {/* <tr>
            <td style={{ backgroundColor:  "#FFE659"}} >{1}</td>
            <td style={{ backgroundColor:  "#FFE659"}} >{1}</td>
            <td style={{ backgroundColor:  "#FFE659"}} >{1}</td>
            <td style={{ backgroundColor:  "#FFE659"}} >{1}</td>
            <td style={{ backgroundColor:  "#FFE659"}} >{1}</td>
            <td style={{ backgroundColor:  "#FFE659"}} >{1}</td>
          </tr> */}
          </tbody>
        </table>
      </div>
      }
    </div>
  );
}

export default AdministradorTareas;
