package targets

import "log"

func init() {

	TargetMap["atsame51-cmsisdap"] = TargetFunc(func(args *Args) error {

		log.Printf("GOT HERE: atsame51-cmsisdap")

		return nil
	})

}
