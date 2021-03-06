package model

import (
    "appengine"
    "appengine/datastore"
	"appengine/memcache"
	"encoding/json"
	//"strconv"
	"strings"
	"net/http"
	"time"
	//"sharded_counter"
)

const GMT = 6
var GMTADJ = -1*3600*GMT

func init() {
    http.HandleFunc("/", home)
}

func home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index.html", http.StatusFound)
}

type Cta struct {
	Folio			int32
	Nombre			string
	Apellidos		string
	Puesto			string
	Email			string
	EmailAlt		string
	Pass			string
	Tel				string
	Cel				string
	FechaHora		time.Time
	UsuarioInt		string
	CodigoCfm		string
	Status			bool
}

type Empresa struct {
	IdEmp		string
	Folio		int32
	RFC			string
	Nombre		string
	RazonSoc	string
	DirCalle	string
	DirCol		string
	DirEnt		string
	DirMun		string
	DirCp		string
	NumSuc		string
	OrgEmp		string
	OrgEmpOtro	string
	OrgEmpReg	string
	//Entidades	[]Entidad
	Url			string
	Benef		int
	PartLinea	int
	ExpComer	int
	Desc		string
	FechaHora	time.Time
	Status		bool
}

type CtaEmpresa struct {
	IdEmp		string
	Email			string
	EmailAlt		string
}

type EmpresaNm struct {
	IdEmp		string
	Folio		int32
	RFC			string
	Nombre		string
	RazonSoc	string
}

type Sucursal struct {
	IdSuc		string
	IdEmp		string
	Nombre		string
	Tel			string
	DirCalle	string
	DirCol		string
	DirEnt		string
	DirMun		string
	DirCp		string
	GeoUrl		string
	Geo1		string
	Geo2		string
	Geo3		string
	Geo4		string
	FechaHora	time.Time
	Latitud		float64
	Longitud	float64
}

type Quest struct {
	PartLinea	int
	ExpComer	int
	Desc		string
}

type Entidad struct {
	CveEnt		string
	Entidad		string
	Abrv		string
	CveCap		string
	Capital		string
	Selected	string
}

type Municipio struct {
	CveEnt		string
	Entidad		string
	Abrv		string
	CveMun		string
	Municipio	string
	CvaCab		string
	Cabecera	string
	Selected	string
}

type Organismo struct {
	Siglas		string
	Nombre		string
	Selected	string
}

type Image struct {
	Data	[]byte
	IdEmp	string
	IdImg	string
	Kind	string
	Name	string
	Desc	string
	Sizepx	int
	Sizepy	int
	Url		string
	Type	string
	Sp1		string
	Sp2		string
	Sp3		string
	Sp4		string
	Np1		int
	Np2		int
	Np3		int
	Np4		int
}

/*
 * Métodos de control de cambios
 */
type ChangeControl struct {
	Id		string
	Kind	string
	Status	string
	FechaHora	time.Time
}

func PutChangeControl(c appengine.Context, id string, kind string, status string) error {
	var cc ChangeControl
	cc.Id = id
	cc.Kind = kind
	cc.Status = status
	cc.FechaHora = time.Now().Add(time.Duration(GMTADJ)*time.Second)
	_, err := datastore.Put(c, datastore.NewKey(c, "ChangeControl", kind+"_"+id, 0, nil), &cc)
	if err != nil {
		return err
	}
	return nil
}

func PutCtaEmp(c appengine.Context, idemp string, email string, emailalt string) error {
	var ce CtaEmpresa
	ce.IdEmp = idemp
	ce.Email = email
	ce.EmailAlt = emailalt
	_, err := datastore.Put(c, datastore.NewKey(c, "CtaEmpresa", ce.IdEmp, 0, nil), &ce)
	if err != nil {
		c.Errorf("PutCtaEmpresa(); Error al intentar crear CtaEmpresa : %v", idemp)
		return err
	}
	return nil
}

func DelCtaEmp(c appengine.Context, idemp string) error {
    if err := datastore.Delete(c, datastore.NewKey(c, "CtaEmpresa", idemp, 0, nil)); err != nil {
		c.Errorf("DelCtaEmpresa(); Error al intentar borrar CtaEmpresa : %v", idemp)
		return err
	}
	return nil
}

/*
 * Métodos de acceso, modificación y limpieza 
 */
func (r *Cta) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Cta", r.Email, 0, nil)
}

func GetCta(c appengine.Context, email string) (*Cta, error) {
	ua := &Cta{ Email: email, }
	err := datastore.Get(c, ua.Key(c), ua)
	if err == datastore.ErrNoSuchEntity {
		return ua, err
	}
	return ua, nil
}

func (r *Cta) DelCta(c appengine.Context) error {
    if err := datastore.Delete(c, r.Key(c)); err != nil {
		return err
	}
	return nil
}

func PutCta(c appengine.Context, u *Cta) (*Cta, error) {
	_, err := datastore.Put(c, u.Key(c), u)
	if err != nil {
		return nil, err
	}
	return u, err
}

func (r *Cta) GetEmpresa(c appengine.Context, id string) (*Empresa, error) {
	e := &Empresa{ IdEmp: id }
	err := datastore.Get(c, datastore.NewKey(c, "Empresa", e.IdEmp, 0, r.Key(c)), e)
	if err == datastore.ErrNoSuchEntity {
		return nil, err
	}
	return e, nil
}

func GetEmpresaM(c appengine.Context, id string) *Empresa {
	e := &Empresa{ IdEmp: id }
	err := datastore.Get(c, datastore.NewKey(c, "EmpLogo", e.IdEmp, 0, nil), e)
	if err == datastore.ErrNoSuchEntity {
		return nil
	}
	return e
}

func (r *Cta) PutEmpresa(c appengine.Context, e *Empresa) (*Empresa, error) {
	//if e.Folio == 0 {
	//	if err := sharded_counter.Increment(c, "empresa"); err == nil {
	//		if olio, err := sharded_counter.Count(c, "empresa"); err == nil {
	//			e.Folio = folio
	//		}
	//	}
	//}

	_, err := datastore.Put(c, datastore.NewKey(c, "Empresa", e.IdEmp, 0, r.Key(c)), e)
	if err != nil {
		return nil, err
	}
	_ = PutChangeControl(c, e.IdEmp, "Empresa", "M")

	/*
	 * Se consulta la empresa normalizada para actualizar datos
	 */
	en := &EmpresaNm{ IdEmp: e.IdEmp }
	err = datastore.Get(c, datastore.NewKey(c, "EmpresaNm", en.IdEmp, 0, r.Key(c)), en)
	if err == datastore.ErrNoSuchEntity {
		// No existe, se crea el registro normalizado
		// Todo esto se hizo porque no se planeo bien desde un principio y 
		// se requiere normalizar el nombre de empresa así como un folio
		// para llevar la cuenta en tiempo real, entre otras cosas :S
		en.Folio = e.Folio
		en.RFC = strings.ToUpper(e.RFC)
		en.Nombre = strings.ToLower(e.Nombre)
		en.RazonSoc = strings.ToLower(e.RazonSoc)
		_, err = datastore.Put(c, datastore.NewKey(c, "EmpresaNm", en.IdEmp, 0, r.Key(c)), en)
		if err != nil {
			c.Errorf("PutEmpresa() Error al intentar crear EmpresaNm : %v", e.IdEmp)
		}
	} else {
		en.RFC = strings.ToUpper(e.RFC)
		en.Nombre = strings.ToLower(e.Nombre)
		en.RazonSoc = strings.ToLower(e.RazonSoc)
		_, err = datastore.Put(c, datastore.NewKey(c, "EmpresaNm", en.IdEmp, 0, r.Key(c)), en)
		if err != nil {
			c.Errorf("PutEmpresa() Error al intentar actualizar EmpresaNm : %v", e.IdEmp)
		}
	}
	return e, nil
}

func (r *Cta) NewEmpresa(c appengine.Context, e *Empresa) (*Empresa, error) {
	/*
		OJO:
		Agregar consulta del random antes de crear entity
	*/
	//if err := sharded_counter.Increment(c, "empresa"); err == nil {
	//	if folio, err := sharded_counter.Count(c, "empresa"); err == nil {
			e.IdEmp = RandId(20)
			//e.Folio = folio
			_, err := datastore.Put(c, datastore.NewKey(c, "Empresa", e.IdEmp, 0, r.Key(c)), e)
			if err != nil {
				return nil, err
			}
			_ = PutChangeControl(c, e.IdEmp, "Empresa", "A")
			c.Infof("Empresa creada Folio: %d, RandID: %v", e.Folio, e.IdEmp)

			var en EmpresaNm
			en.IdEmp = e.IdEmp
			en.Folio = e.Folio
			en.RFC = strings.ToUpper(e.RFC)
			en.Nombre = strings.ToLower(e.Nombre)
			en.RazonSoc = strings.ToLower(e.RazonSoc)
			_, err = datastore.Put(c, datastore.NewKey(c, "EmpresaNm", en.IdEmp, 0, r.Key(c)), &en)
			if err != nil {
				c.Errorf("Error al intentar crear EmpresaNm : %v", e.IdEmp)
				return nil, err
			}
	//	} else {
	//		c.Errorf("Folio inseguro al intentar crear Empresa : %v", e.IdEmp)
	//		return e, err
	//	}
	//} else {
//		c.Errorf("Folio inseguro al intentar crear Empresa : %v", e.IdEmp)
//		return e, err
//	}
	return e, nil
}

func (r *Cta) DelEmpresa(c appengine.Context, id string) error {
	if err := DelImg(c, "EmpLogo", id); err != nil {
		return err
	}
	if err := DelImg(c, "ShortLogo", id); err != nil {
		return err
	}
	if err := DelSucs(c, id); err != nil {
		return err
	}
    if err := datastore.Delete(c, datastore.NewKey(c, "Empresa", id, 0, r.Key(c))); err != nil {
		return err
	}
    if err := datastore.Delete(c, datastore.NewKey(c, "EmpresaNm", id, 0, r.Key(c))); err != nil {
		return err
	}
    if err := datastore.Delete(c, datastore.NewKey(c, "CtaEmpresa", id, 0, nil)); err != nil {
		c.Errorf("DelCtaEmpresa(); Error al intentar borrar CtaEmpresa : %v", id)
		return err
	}
	_ = PutChangeControl(c, id, "Empresa", "B")
	return nil
}

// Métodos de Empresa
func GetEmpresa(c appengine.Context, id string) (*Empresa) {
	/* llave de Cta-Empresa */
	ce := &CtaEmpresa{ IdEmp: id }
	ceKey := datastore.NewKey(c, "CtaEmpresa", ce.IdEmp, 0, nil)
	err := datastore.Get(c, ceKey, ce)
	if err == datastore.ErrNoSuchEntity {
		return nil
	}

	/* parent de Empresa */
	cta := &Cta{ Email: ce.Email }
	ctaKey := datastore.NewKey(c, "Cta", cta.Email, 0, nil)
	err = datastore.Get(c, ctaKey, cta)
	if err == datastore.ErrNoSuchEntity {
		return nil
	}

	/* Key de empresa */
	emp := &Empresa{ IdEmp: id }
	empKey := datastore.NewKey(c, "Empresa", emp.IdEmp, 0, ctaKey)
	err = datastore.Get(c, empKey, emp)
	if err == datastore.ErrNoSuchEntity {
		return nil
	}
	return emp
}

func GetEmpSucursales(c appengine.Context, IdEmp string) *[]Sucursal {
	/* llave de Cta-Empresa */
	ce := &CtaEmpresa{ IdEmp: IdEmp }
	ceKey := datastore.NewKey(c, "CtaEmpresa", ce.IdEmp, 0, nil)
	err := datastore.Get(c, ceKey, ce)
	if err == datastore.ErrNoSuchEntity {
		return nil
	}

	/* parent de Empresa */
	cta := &Cta{ Email: ce.Email }
	ctaKey := datastore.NewKey(c, "Cta", cta.Email, 0, nil)
	err = datastore.Get(c, ctaKey, cta)
	if err == datastore.ErrNoSuchEntity {
		return nil
	}

	/* Key de empresa */
	empKey := datastore.NewKey(c, "Empresa", IdEmp, 0, ctaKey)
	q := datastore.NewQuery("Sucursal").Ancestor(empKey)
	n, _ := q.Count(c)
	sucursales := make([]Sucursal, 0, n)
	if _, err := q.GetAll(c, &sucursales); err != nil {
		return nil
	}
	return &sucursales
}

func TouchSuc(c appengine.Context, IdSuc string, IdEmp string) error {
	/* llave de Cta-Empresa */
	ce := &CtaEmpresa{ IdEmp: IdEmp }
	ceKey := datastore.NewKey(c, "CtaEmpresa", ce.IdEmp, 0, nil)
	err := datastore.Get(c, ceKey, ce)
	if err == datastore.ErrNoSuchEntity {
		return nil
	}

	/* parent de Empresa */
	cta := &Cta{ Email: ce.Email }
	ctaKey := datastore.NewKey(c, "Cta", cta.Email, 0, nil)
	err = datastore.Get(c, ctaKey, cta)
	if err == datastore.ErrNoSuchEntity {
		return nil
	}

	/* Key de empresa */
	empKey := datastore.NewKey(c, "Empresa", IdEmp, 0, ctaKey)

	suc := &Sucursal{ IdSuc: IdSuc }
	sucKey := datastore.NewKey(c, "Sucursal", IdSuc, 0, empKey)
	err = datastore.Get(c, sucKey, suc)
	if err == datastore.ErrNoSuchEntity {
		return err
	}
	suc.FechaHora = time.Now().Add(time.Duration(GMTADJ)*time.Second)
	if _, err = datastore.Put(c, sucKey, suc); err != nil {
		return err
	}
	return nil
}

func (e *Empresa) PutSuc(c appengine.Context, cta *Cta, s *Sucursal, idemp string) (*Sucursal, error) {
	if(s.IdSuc == "") {
		s.IdSuc = RandId(20)
		_ = PutChangeControl(c, s.IdSuc, "Sucursal", "A")
	} else {
		_ = PutChangeControl(c, s.IdSuc, "Sucursal", "M")
	}
	empKey := datastore.NewKey(c, "Empresa", idemp, 0, cta.Key(c))
    _, err := datastore.Put(c, datastore.NewKey(c, "Sucursal", s.IdSuc, 0, empKey), s)
	if err != nil {
		return nil, err
	}
	return s, err
}

// Métodos de Sucursal
func GetSuc(c appengine.Context, cta *Cta, idsuc string, idemp string) (*Sucursal) {
	suc := &Sucursal{ IdSuc: idsuc }
	sucKey := datastore.NewKey(c, "Sucursal", idsuc, 0, datastore.NewKey(c, "Empresa", idemp, 0, cta.Key(c)))
	err := datastore.Get(c, sucKey, suc)
	if err == datastore.ErrNoSuchEntity {
		// Regresa un cascarón
		return nil
	}
	return suc
}

func DelSuc(c appengine.Context, id string) error {
	q := datastore.NewQuery("Sucursal").Filter("IdSuc =", id)
	for i := q.Run(c); ; {
		var e Sucursal
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		if err:= DelSucursalesOferta(c, e.IdSuc); err != nil {
			return err
		}
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
		_ = PutChangeControl(c, e.IdSuc, "Sucursal", "B")
	}
	return nil
}

func DelSucs(c appengine.Context, idEmp string) error {
	q := datastore.NewQuery("Sucursal").Filter("IdEmp =", idEmp)
	for i := q.Run(c); ; {
		var e Sucursal
		key, err := i.Next(&e)
		if err == datastore.Done {
			break
		}
		if err:= DelSucursalesOferta(c, e.IdSuc); err != nil {
			return err
		}
		if err := datastore.Delete(c, key); err != nil {
			return err
		}
		_ = PutChangeControl(c, e.IdSuc, "Sucursal", "B")
	}
	return nil
}

// Métodos de Entidad
func (e *Entidad) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Entidad", e.CveEnt, 0, nil)
}

func (e *Entidad) GetMunicipios(c appengine.Context) (*[]Municipio, error) {
	q := datastore.NewQuery("Municipio").Ancestor(e.Key(c))
	nm, _ := q.Count(c)
	municipios := make([]Municipio, 0, nm)
	if _, err := q.GetAll(c, &municipios); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, err
		}
	}
	return &municipios, nil
}

func GetEntidad(c appengine.Context, cveent string) (*Entidad, error) {
	e := &Entidad{ CveEnt: cveent }
	err := datastore.Get(c, e.Key(c), e)
	if err == datastore.ErrNoSuchEntity {
		return nil, err
	}
	return e, nil
}



// Métodos de Municipio
func (m *Municipio) Parent(c appengine.Context) *Entidad {
	e, _ := GetEntidad(c, m.CveEnt)
	return e
}

func (m *Municipio) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Municipio", m.CveMun, 0, m.Parent(c).Key(c))
}

func GetMunicipio(c appengine.Context, cvemun string) *Municipio {
	e := &Municipio{ CveMun: cvemun }
	err := datastore.Get(c, e.Key(c), e)
	if err == datastore.ErrNoSuchEntity {
		return nil
	}
	return e
}

// Métodos de Imagen
// Obtiene la llave de una imagen
func (i *Image) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, i.Kind, i.IdImg, 0, nil)
}

// Borra una imagen
func DelImg(c appengine.Context, kind string, id string) error {
    if err := datastore.Delete(c, datastore.NewKey(c, kind, id, 0, nil)); err != nil {
		return err
	}
	return nil
}

// Guarda Imagen modificada
func PutLogo(c appengine.Context, i *Image) (*datastore.Key, error) {
	key, err := datastore.Put(c, datastore.NewKey(c, i.Kind, i.IdEmp, 0, nil), i)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func GetLogo(c appengine.Context, idemp string) (*Image) {
	i := &Image{ IdEmp: idemp, Kind: "EmpLogo" }
	// Para el logo sólo se utiliza la llave IdEmp
	err := datastore.Get(c, datastore.NewKey(c, i.Kind, i.IdEmp, 0, nil), i)
	if err == datastore.ErrNoSuchEntity {
		return nil
	}
	return i
}

func GetShortLogo(c appengine.Context, idemp string) (*Image, error) {
	i := &Image{ IdEmp: idemp, Kind: "ShortLogo" }
	// Para el logo sólo se utiliza la llave IdEmp
	err := datastore.Get(c, datastore.NewKey(c, i.Kind, i.IdEmp, 0, nil), i)
	if err == datastore.ErrNoSuchEntity {
		return nil, err
	}
	return i, err
}

// Obtiene una imagen
func GetImg(c appengine.Context, id string) (*Image, error) {
	i := &Image{ IdImg: id }
	err := datastore.Get(c, i.Key(c), i)
	if err == datastore.ErrNoSuchEntity {
		//_, err = datastore.Put(c, ua.Key(), ua)
		return i, err
	}
	return i, err
}

// Lista entidades
func ListEnt(c appengine.Context, ent string) *[]Entidad {
	estados := make([]Entidad, 0, 32)
	if item, err := memcache.Get(c, "estados"); err == memcache.ErrCacheMiss {
		q := datastore.NewQuery("Entidad")
		if _, err := q.GetAll(c, &estados); err != nil {
			return nil
		}
		b, _ := json.Marshal(estados)
		item := &memcache.Item{
			Key:   "estados",
			Value: b,
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Errorf("memcache.Add Entidad : %v", err)
		}
	} else {
		//c.Infof("Memcache activo: %v", item.Key)
		if err := json.Unmarshal(item.Value, &estados); err != nil {
			c.Errorf("Memcache Unmarshalling item: %v", err)
		}
	}
	for i, _ := range estados {
		if(ent == estados[i].CveEnt) {
			estados[i].Selected = `selected`
		}
	}
	return &estados
}


