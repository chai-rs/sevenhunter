module.exports = {
  /**
   * @param db {import('mongodb').Db}
   * @param client {import('mongodb').MongoClient}
   * @returns {Promise<void>}
   */
  async up(db, client) {
    await db.createCollection('users', {
      validator: {
        $jsonSchema: {
          bsonType: 'object',
          required: ['_id', 'email', 'name', 'hashed_password', 'created_at'],
          properties: {
            _id: {
              bsonType: 'objectId',
              description: 'auto-generated unique identifier'
            },
            email: {
              bsonType: 'string',
              description: 'must be a string and is required'
            },
            hashed_password: {
              bsonType: 'string',
              description: 'must be a string and is required'
            },
            name: {
              bsonType: 'string',
              description: 'must be a string if provided'
            },
            created_at: {
              bsonType: 'date',
              description: 'must be a date and is required'
            },
          }
        }
      }
    });

    // Create unique index on email
    await db.collection('users').createIndex({ email: 1 }, { unique: true });

    // Create index on name if it exists
    await db.collection('users').createIndex({ name: 1 }, { sparse: true });
  },

  /**
   * @param db {import('mongodb').Db}
   * @param client {import('mongodb').MongoClient}
   * @returns {Promise<void>}
   */
  async down(db, client) {
    await db.collection('users').drop();
  }
};
