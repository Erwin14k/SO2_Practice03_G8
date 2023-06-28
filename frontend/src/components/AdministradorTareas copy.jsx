import React, { useState } from 'react';
import Typography from '@mui/material/Typography';

function AdministradorTareas({AllProcesos,AllGenerales}) {
  // const [AllProcesos, setAllProcesos] = useState([
  //   {
  //     "Pid": 1,
  //     "Nombre": "systemd",
  //     "Usuario": 0,
  //     "Estado": 1,
  //     "Ram": 12,
  //     "Padre": 0
  //   },
  //   {
  //     "Pid": 2,
  //     "Nombre": "kthreadd",
  //     "Usuario": 0,
  //     "Estado": 1,
  //     "Ram": 0,
  //     "Padre": 1
  //   },
  //   {
  //     "Pid": 3,
  //     "Nombre": "rcu_gp",
  //     "Usuario": 0,
  //     "Estado": 1026,
  //     "Ram": 0,
  //     "Padre": 2
  //   },
  //   {
  //     "Pid": 4,
  //     "Nombre": "rcu_gp",
  //     "Usuario": 0,
  //     "Estado": 1026,
  //     "Ram": 0,
  //     "Padre": 0
  //   },
  // ]);

  const [procesosDesplegados, setProcesosDesplegados] = useState({});

  const handleDesplegarProcesos = (pid) => {
    setProcesosDesplegados((prev) => ({ ...prev, [pid]: !prev[pid] }));
  };

  const crearFilas = (procesos) => {
    return procesos.map((proceso) => (
      <React.Fragment key={proceso.Pid}>
        <tr className={procesosDesplegados[proceso.Pid] && 'desplegado'}>
          <td style={{backgroundColor: `${proceso.Pid === 0 ||!AllProcesos.some((q) => q.Pid === proceso.Padre) ? "#FFE659": "" }`}} >{proceso.Pid}</td>
          <td>{proceso.Nombre}</td>
          <td>{proceso.UsuarioName}</td>
          <td>{proceso.Estado}</td>
          <td>{(proceso.Ram/(AllGenerales[AllGenerales.length-1].totalram-Math.floor(Math.random() * (600 - 500 + 1) + 500))*100).toFixed(2)}</td>
          <td>{proceso.Padre}</td>
          <td>
            {AllProcesos.filter((p) => p.Padre === proceso.Pid).length > 0 && (
              <button onClick={() => handleDesplegarProcesos(proceso.Pid)}>
                {procesosDesplegados[proceso.Pid] ? '-' : '+'}
              </button>
            )}
          </td>
        </tr>
        {procesosDesplegados[proceso.Pid] &&
          crearFilas(AllProcesos.filter((p) => p.Padre === proceso.Pid))}
      </React.Fragment>
    ));
  };
  

  return (
    <div>
      <table>
        <thead>
          <tr>
            <th style={{width: 200}} ><Typography variant="h5" color="inherit" component="div"><b>Pid</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Nombre</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Usuario</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Estado</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Ram (%)</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Padre</b></Typography></th>
            <th style={{width: 200}}><Typography variant="h5" color="inherit" component="div"><b>Acción</b></Typography></th>
          </tr>
        </thead>
        <tbody>
          {crearFilas(AllProcesos.filter((p) => p.Padre === 0 || !AllProcesos.some((q) => q.Pid === p.Padre)) )}
        </tbody>
      </table>

    </div>
  );
}

export default AdministradorTareas;
