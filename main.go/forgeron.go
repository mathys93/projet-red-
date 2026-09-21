# Inventaire du joueur
joueur = {
    "Fragment": 10,          # fragments disponibles
    "Inventaire": []         # objets possédés
}

# Objets vendus par le forgeron (prix en fragments)
forgeron_objets = {
    "Chapeau de l'aventurier": 2,
    "Tunique de l'aventurier": 3,
    "Bottes de l'aventurier": 2
}

# Vérifie si le joueur peut acheter
def peut_acheter(item):
    if item not in forgeron_objets:
        print("Le forgeron ne vend pas cet objet.")
        return False
    
    prix = forgeron_objets[item]
    
    if joueur["Fragment"] < prix:
        print(f"Tu n'as pas assez de fragments pour acheter {item}.")
        print(f"Il faut {prix} fragments, tu n'en as que {joueur['Fragment']}.")
        return False
    
    return True

# Achat de l'objet
def acheter(item):
    if peut_acheter(item):
        prix = forgeron_objets[item]
        joueur["Fragment"] -= prix
        joueur["Inventaire"].append(item)
        
        print(f"Tu as acheté {item} !")
        print(f"Fragments restants : {joueur['Fragment']}")
        print(f"Inventaire du joueur : {joueur['Inventaire']}")
    else:
        print("Achat impossible.")

# --- EXEMPLES D'ACHAT ---
acheter("Chapeau de l'aventurier")
acheter("Tunique de l'aventurier")
acheter("Bottes de l'aventurier")