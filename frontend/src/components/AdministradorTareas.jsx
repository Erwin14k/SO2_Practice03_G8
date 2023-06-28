import React, { useState } from 'react';
import Typography from '@mui/material/Typography';
import { CanvasJSChart } from 'canvasjs-react-charts'
import Paper from '@mui/material/Paper';

import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import Button from '@mui/material/Button';

function AdministradorTareas({ AllProcesos, AllGenerales }) {
  const [pid, setPid] = useState(null);
  const [arrMemory, setArrMemory] = useState([]);
  const [procesosDesplegados, setProcesosDesplegados] = useState({});

  const [inicioMem, setInicioMem] = useState('');
  const [finMem, setFinMem] = useState('');

  const [valueRSS, setValueRSS] = useState(0);
  const [valueSize, setValueSize] = useState(0);
  const [valuePercent, setValuePercent] = useState(0);

  const options = {
    animationEnabled: true,
    exportEnabled: true,
    theme: "light2", // "light1", "dark1", "dark2"
    title: {
      text: `Fin: ${finMem}`
    },
    axisY: {
      title: "MEMORIA USADA",
      suffix: "MB"
    },
    axisX: {
      title: `Inicio: ${inicioMem}`,
      prefix: "",
      interval: 0
    },
    data: [
      {
        type: "stackedColumn",
        name: "vmRSS",
        showInLegend: true,
        yValueFormatString: "#,###.###MB",
        dataPoints: [
          { label: "vmRSS", y: valueRSS }
        ]
      },
      {
        type: "stackedColumn",
        name: "vmSize",
        showInLegend: true,
        yValueFormatString: "#,###.###MB",
        dataPoints: [
          { label: "vmSize", y: valueSize }
        ]
      }]
  }

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

  const handleSeeMemory = async (pid,percent) => {
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
        console.log(data);

        setPid(pid);
        setValueRSS(data.total_rss);
        setValueSize(data.total_size);
        setInicioMem("");
        setFinMem("");
        setValuePercent(percent);

        setArrMemory(data.blocks);
        // <td style={{ backgroundColor: `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}`}} >{memoria.initial_address}</td>
        // <td style={{ backgroundColor:  `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}`}} >{memoria.final_address}</td>
        setInicioMem(data.blocks.find(block => block.initial_address !== "").initial_address)
        for (let i = data.blocks.length - 1; i >= 0; i--) {
          if (data.blocks[i].initial_address !== "") {
            setFinMem(data.blocks[i].final_address)
            break;
          }
        }


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
            <button onClick={() => handleSeeMemory(proceso.pid,(proceso.ram / (AllGenerales[AllGenerales.length - 1].totalram - Math.floor(Math.random() * (600 - 500 + 1) + 500)) * 100).toFixed(2))}>
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
                <button style={{ marginRight: "10px", marginBottom: "10px" }} onClick={() => closeMemory()}>
                  cerrar
                </button>
                Asignación de memoria del proceso - PID: {pid}

              </Typography>

            </center>
          </Paper>
          <div className="container2">
            <div className='centerSegMem2'>

              <Card sx={{ minWidth: 275 }}>
                <CardContent>
                  <Typography variant="h5" component="div">
                    Memoria Residente
                  </Typography>
                  <Typography variant="body2">
                    {`${valueRSS.toFixed(3)} MB`}
                  </Typography>
                </CardContent>
              </Card>

              <br />

              <Card sx={{ minWidth: 275 }}>
                <CardContent>
                  <Typography variant="h5" component="div">
                    Memoria Virtual
                  </Typography>
                  <Typography variant="body2">
                    {`${valueSize.toFixed(3)} MB`}
                  </Typography>
                </CardContent>
              </Card>

              <br />
              
              <Card sx={{ minWidth: 275 }}>
                <CardContent>
                  <Typography variant="h5" component="div">
                    % de consumo
                  </Typography>
                  <Typography variant="body2">
                  {`${valuePercent} %`}
                  </Typography>
                </CardContent>
              </Card>
            </div>

            <div className='centerSegMem' >
              <Card sx={{ minWidth: 275 }}>
                <CanvasJSChart options={options} />
              </Card>
            </div>

          </div>



          <table>
            <thead>
              <tr>
                <th style={{ width: 200 }} ><Typography variant="h5" color="inherit" component="div"><b>Dispositivo</b></Typography></th>
                <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Archivo</b></Typography></th>
                <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Direccion de memoria virtua Inicial</b></Typography></th>
                <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Direccion de memoria virtua Final</b></Typography></th>
                <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Permisos</b></Typography></th>
                <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>RSS (MB)</b></Typography></th>
                <th style={{ width: 200 }}><Typography variant="h5" color="inherit" component="div"><b>Tamaño (MB)</b></Typography></th>
              </tr>
            </thead>
            <tbody>
              {
                arrMemory !== null &&
                arrMemory.map((memoria, index) => (
                  <tr key={index}>
                    <td style={{ backgroundColor: `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}` }} >{memoria.device}</td>
                    <td style={{ backgroundColor: `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}` }} >{memoria.file}</td>
                    <td style={{ backgroundColor: `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}` }} >{memoria.initial_address}</td>
                    <td style={{ backgroundColor: `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}` }} >{memoria.final_address}</td>
                    <td style={{ backgroundColor: `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}` }} >{memoria.permissions != null ? memoria.permissions.join(" - ") : []}</td>
                    <td style={{ backgroundColor: `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}` }} >{memoria.rss}</td>
                    <td style={{ backgroundColor: `${index % 2 === 0 ? "#FFEE8C" : "#D8D8D8"}` }} >{memoria.size}</td>
                  </tr>
                ))

              }
            </tbody>
          </table>
        </div>
      }
    </div>
  );
}

export default AdministradorTareas;
