package services

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/gestion_espacios_fisicos_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

func RegistrarEspacioFisico(transaccion *models.NuevoEspacioFisico) (alerta []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"funcion": "RegistrarEspacioFisico", "err": err, "status": "500"}
			panic(outputError)
		}
	}()
	alerta = append(alerta, "Success")

	var creaciones models.Creaciones

	var tipoEspacioFisico models.TipoEspacioFisico
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_espacio_fisico/" + strconv.Itoa(transaccion.TipoEspacioFisico)
	if err := request.GetJson(url, &tipoEspacioFisico); err != nil || tipoEspacioFisico.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	var dependencia models.Dependencia
	url = beego.AppConfig.String("OikosCrudUrl") + "dependencia/" + strconv.Itoa(transaccion.DependenciaPadre)
	if err := request.GetJson(url, &dependencia); err != nil || tipoEspacioFisico.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	var tipoUso models.TipoUso
	url = beego.AppConfig.String("OikosCrudUrl") + "tipo_uso/" + strconv.Itoa(transaccion.TipoUso)
	if err := request.GetJson(url, &tipoUso); err != nil || tipoEspacioFisico.Id == 0 {
		logs.Error(err)
		panic(err.Error())
	}

	var espacioFisico = CrearEspacioFisico(transaccion, tipoEspacioFisico, &creaciones)
	
	CrearAsignacionEspacioFisicoDependencia(transaccion, dependencia, espacioFisico, &creaciones)
	CrearTipoUsoEspacioFisico(transaccion, tipoUso, espacioFisico, &creaciones)
	fmt.Println("ESPACIO FISICO ID")
	fmt.Println(creaciones.EspacioFisicoId)
	fmt.Println(espacioFisico)
	fmt.Println("ASIGNACION ESPACIO FISICO DEPENDENCIA ID")
	fmt.Println(creaciones.AsignacionEspacioFisicoDependenciaId)
	fmt.Println("TIPO USO FISICO DEPENDENCIA ID")
	fmt.Println(creaciones.TipoUsoEspacioFisico)

	if len(transaccion.CamposDinamicos) != 0 {
		var camposCreados = CrearCampos(transaccion, &creaciones)
		var espacioFisicoCampoCreados = CrearEspacioFisicoCampo(transaccion, camposCreados, espacioFisico, &creaciones)
		fmt.Println("CAMPOOOOOS")
		fmt.Println(camposCreados)
		fmt.Println(creaciones.CamposId)
		fmt.Println("ESPACIOS FISICOS CAMPOOOOOS")
		fmt.Println(espacioFisicoCampoCreados)
		fmt.Println(creaciones.EspacioFisicoCampoId)
	}

	return alerta, outputError
}

func CrearEspacioFisico(transaccion *models.NuevoEspacioFisico, tipoEspacioFisico models.TipoEspacioFisico, creaciones *models.Creaciones) (nuevoEspacioFisico models.EspacioFisico) {
	nuevoEspacioFisico.Nombre = transaccion.EspacioFisico.Nombre
	nuevoEspacioFisico.Descripcion = transaccion.EspacioFisico.Descripcion
	nuevoEspacioFisico.CodigoAbreviacion = transaccion.EspacioFisico.CodigoAbreviacion
	nuevoEspacioFisico.TipoTerrenoId = transaccion.TipoTerreno
	nuevoEspacioFisico.TipoEdificacionId = transaccion.TipoEdificacion
	nuevoEspacioFisico.TipoEspacioFisicoId = &tipoEspacioFisico
	nuevoEspacioFisico.Activo = true
	nuevoEspacioFisico.FechaCreacion = time_bogota.TiempoBogotaFormato()
	nuevoEspacioFisico.FechaModificacion = time_bogota.TiempoBogotaFormato()
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico"
	var resEspacioFisicoRegistrado map[string]interface{}
	if err := request.SendJson(url, "POST", &resEspacioFisicoRegistrado, nuevoEspacioFisico); err != nil || resEspacioFisicoRegistrado["Id"] == nil{
		logs.Error(err)
		panic(err.Error())
	}
	creaciones.EspacioFisicoId = int(resEspacioFisicoRegistrado["Id"].(float64))
	nuevoEspacioFisico.Id = int(resEspacioFisicoRegistrado["Id"].(float64))
	fmt.Println("ESPACIO FISICO CREADO")
	return nuevoEspacioFisico
}

func CrearAsignacionEspacioFisicoDependencia(transaccion *models.NuevoEspacioFisico, dependencia models.Dependencia, espacioFisico models.EspacioFisico, creaciones *models.Creaciones) {
	var asignacionEspacioFisicoDependencia models.AsignacionEspacioFisicoDependencia
	asignacionEspacioFisicoDependencia.EspacioFisicoId = &espacioFisico
	asignacionEspacioFisicoDependencia.DependenciaId = &dependencia
	asignacionEspacioFisicoDependencia.DocumentoSoporte = 0
	asignacionEspacioFisicoDependencia.FechaInicio = time_bogota.TiempoBogotaFormato()
	asignacionEspacioFisicoDependencia.FechaFin = time_bogota.TiempoBogotaFormato()
	asignacionEspacioFisicoDependencia.Activo = true
	asignacionEspacioFisicoDependencia.FechaCreacion = time_bogota.TiempoBogotaFormato()
	asignacionEspacioFisicoDependencia.FechaModificacion = time_bogota.TiempoBogotaFormato()
	url := beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia"
	var resAsignacionEspacioFisicoDependenciaRegistrado map[string]interface{}
	if err := request.SendJson(url, "POST", &resAsignacionEspacioFisicoDependenciaRegistrado, asignacionEspacioFisicoDependencia); err != nil || resAsignacionEspacioFisicoDependenciaRegistrado["Id"] == nil{
		fmt.Println("ENTRA A ERROR DE ASIGNACION")
		rollbackEspacioFisicoCreado(creaciones)
		logs.Error(err)
		panic(err.Error())
	}
	fmt.Println(resAsignacionEspacioFisicoDependenciaRegistrado["Id"])
	fmt.Println("ASIGNACION CREADO")
	creaciones.AsignacionEspacioFisicoDependenciaId = int(resAsignacionEspacioFisicoDependenciaRegistrado["Id"].(float64))
}

func CrearTipoUsoEspacioFisico(transaccion *models.NuevoEspacioFisico, tipoUso models.TipoUso, espacioFisico models.EspacioFisico, creaciones *models.Creaciones) {
	var tipoUsoEspacioFisico models.TipoUsoEspacioFisico
	tipoUsoEspacioFisico.EspacioFisicoId = &espacioFisico
	tipoUsoEspacioFisico.TipoUsoId = &tipoUso
	tipoUsoEspacioFisico.Activo = true
	tipoUsoEspacioFisico.FechaCreacion = time_bogota.TiempoBogotaFormato()
	tipoUsoEspacioFisico.FechaModificacion = time_bogota.TiempoBogotaFormato()
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico"
	var resTipoUsoEspacioFisicoRegistrado map[string]interface{}
	if err := request.SendJson(url, "POST", &resTipoUsoEspacioFisicoRegistrado, tipoUsoEspacioFisico); err != nil || resTipoUsoEspacioFisicoRegistrado["Id"] == nil{
		rollbackAsignacionEspacioFisicoDependencia(creaciones)
		logs.Error(err)
		panic(err.Error())
	}
	fmt.Println(resTipoUsoEspacioFisicoRegistrado["Id"])
	fmt.Println("TIPO USO CREADO")
	creaciones.TipoUsoEspacioFisico = int(resTipoUsoEspacioFisicoRegistrado["Id"].(float64))
}

func CrearCampos(transaccion *models.NuevoEspacioFisico, creaciones *models.Creaciones) (campos []models.Campo) {
	for _, campo := range transaccion.CamposDinamicos {
		var campoDinamico models.Campo
		campoDinamico.Nombre = campo.NombreCampo
		campoDinamico.Descripcion = campo.Descripcion
		campoDinamico.CodigoAbreviacion = campo.CodigoAbreviacion
		campoDinamico.Activo = true
		campoDinamico.FechaCreacion = time_bogota.TiempoBogotaFormato()
		campoDinamico.FechaModificacion = time_bogota.TiempoBogotaFormato()
		url := beego.AppConfig.String("OikosCrudUrl") + "campo"
		var resCampoRegistrado map[string]interface{}
		if err := request.SendJson(url, "POST", &resCampoRegistrado, campoDinamico); err != nil || resCampoRegistrado["Id"] == nil{
			if len(creaciones.CamposId) == 0 {
				rollbackTipoUsoEspacioFisico(creaciones)
			} else {
				fmt.Println("ENTRA A A ERROR DE CAMPOOOS")
				rollbackCrearCampos(creaciones)
			}
			logs.Error(err)
			panic(err.Error())
		}
		fmt.Println(resCampoRegistrado["Id"])
		fmt.Println("CAMPO CREADO")
		creaciones.CamposId = append(creaciones.CamposId, int(resCampoRegistrado["Id"].(float64)))
		campoDinamico.Id = int(resCampoRegistrado["Id"].(float64))
		campos = append(campos, campoDinamico)

	}
	return campos
}

func CrearEspacioFisicoCampo(transaccion *models.NuevoEspacioFisico, camposEntrada []models.Campo, espacioFisico models.EspacioFisico, creaciones *models.Creaciones) (campos []models.EspacioFisicoCampo) {
	for _, campo := range camposEntrada {
		var espacioFisicoCampo models.EspacioFisicoCampo
		espacioFisicoCampo.Valor = "0"
		espacioFisicoCampo.EspacioFisicoId = &espacioFisico
		espacioFisicoCampo.CampoId = &campo
		espacioFisicoCampo.Activo = true
		espacioFisicoCampo.FechaInicio = time_bogota.TiempoBogotaFormato()
		espacioFisicoCampo.FechaFin = time_bogota.TiempoBogotaFormato()
		espacioFisicoCampo.FechaCreacion = time_bogota.TiempoBogotaFormato()
		espacioFisicoCampo.FechaModificacion = time_bogota.TiempoBogotaFormato()
		url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo"
		var resEspacioFisicoCampoRegistrado map[string]interface{}
		if err := request.SendJson(url, "POST", &resEspacioFisicoCampoRegistrado, espacioFisicoCampo); err != nil || resEspacioFisicoCampoRegistrado["Id"] == nil{
			if len(creaciones.EspacioFisicoCampoId) == 0 {
				rollbackCrearCampos(creaciones)
			} else {
				rollbackEspacioFisicoCampo(creaciones)
			}
			logs.Error(err)
			panic(err.Error())
		}
		fmt.Println(resEspacioFisicoCampoRegistrado["Id"])
		fmt.Println("ESPACIO FISICO CAMPO CREADO")
		creaciones.EspacioFisicoCampoId = append(creaciones.EspacioFisicoCampoId, int(resEspacioFisicoCampoRegistrado["Id"].(float64)))
		espacioFisicoCampo.Id = int(resEspacioFisicoCampoRegistrado["Id"].(float64))
		campos = append(campos, espacioFisicoCampo)
	}
	return campos
}

func rollbackEspacioFisicoCreado(creaciones *models.Creaciones) (outputError map[string]interface{}) {
	var respuesta map[string]interface{}
	url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico/" + strconv.Itoa(creaciones.EspacioFisicoId)
	if err := request.SendJson(url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback del espacio fisico" + err.Error())
	}
	return nil
}

func rollbackAsignacionEspacioFisicoDependencia(creaciones *models.Creaciones) (outputError map[string]interface{}) {
	var respuesta map[string]interface{}
	url := beego.AppConfig.String("OikosCrudUrl") + "asignacion_espacio_fisico_dependencia/" + strconv.Itoa(creaciones.AsignacionEspacioFisicoDependenciaId)
	if err := request.SendJson(url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback de la asignacion del espacio fisico dependencia" + err.Error())
	}
	rollbackEspacioFisicoCreado(creaciones)
	return nil
}

func rollbackTipoUsoEspacioFisico(creaciones *models.Creaciones) (outputError map[string]interface{}) {
	var respuesta map[string]interface{}
	url := beego.AppConfig.String("OikosCrudUrl") + "tipo_uso_espacio_fisico/" + strconv.Itoa(creaciones.TipoUsoEspacioFisico)
	if err := request.SendJson(url, "DELETE", &respuesta, nil); err != nil {
		panic("Rollback del tipo de uso del espacio fisico" + err.Error())
	}
	rollbackAsignacionEspacioFisicoDependencia(creaciones)
	return nil
}

func rollbackCrearCampos(creaciones *models.Creaciones) (outputError map[string]interface{}) {
	for _, campo := range creaciones.CamposId {
		var respuesta map[string]interface{}
		url := beego.AppConfig.String("OikosCrudUrl") + "campo/" + strconv.Itoa(campo)
		if err := request.SendJson(url, "DELETE", &respuesta, nil); err != nil {
			panic("Rollback del crear campo" + err.Error())
		}
	}
	rollbackTipoUsoEspacioFisico(creaciones)
	return nil
}

func rollbackEspacioFisicoCampo(creaciones *models.Creaciones) (outputError map[string]interface{}) {
	for _, campo := range creaciones.EspacioFisicoCampoId {
		var respuesta map[string]interface{}
		url := beego.AppConfig.String("OikosCrudUrl") + "espacio_fisico_campo/" + strconv.Itoa(campo)
		if err := request.SendJson(url, "DELETE", &respuesta, nil); err != nil {
			panic("Rollback del espacio fisico campo" + err.Error())
		}
	}
	rollbackCrearCampos(creaciones)
	return nil
}
